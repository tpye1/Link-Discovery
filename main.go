package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type connect_struct struct {
	id           int
	name         string
	network_card string
	mac_addr     string
	status       bool
	ip_addr      string
	iface        *net.Interface
}

type valid_iface struct {
	is_running bool
	iface      *net.Interface
}


type LinkDataHelperForPackets struct {
	SwitchName    string `json:"switch_name"`
	PortID        string `json:"port_id"`
	VlanID        string `json:"vlan_id"`
	SwitchIP      string `json:"switch_ip"`
	SwitchModel   string `json:"switch_model"`
	DuplexOption  string `json:"duplex_option"`
	VptMgmtDomain string `json:"vpt_mgmt_domain"`
	Protocol      string `json:"protocol"`
}

type link_data struct {
	connection      *connect_struct
	switch_name     string
	port_id         string
	vlan_id         string
	switch_ip       string
	switch_model    string
	duplex_option   string
	vpt_mgmt_domain string
	protocol        string
	//stacked			bool
	//member			int
}


func main() {

	app:= tview.NewApplication()
	
	connections := get_connection_data()
	
	var connect_map = make(map[string]*connect_struct)

	var connect_names_arr []string
	
	for i:=0; i < len(connections); i++ {
		connect_map[connections[i].name] = &(connections[i])
		connect_names_arr = append(connect_names_arr, connections[i].name)
	}

	form := tview.NewForm()

	var result_found bool = false 

	var result link_data

	form.SetBorder(true).SetTitle("Link Discovery Client for Linux").SetTitleAlign(tview.AlignCenter)

	// Connection options
	form.AddDropDown("Connection: ", connect_names_arr, 0, nil)

	// Save Link data button
	form.AddButton("Save Link data", func() {
		if result_found == false {return}
		save_link_data(&result)

	})

	// Get Link data button
	form.AddButton("Get Link Data", func() {
		_, selected_str := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
		connection_pointer := connect_map[selected_str]
		result, result_found = get_link_data(connection_pointer)
	})

	// Help button
	form.AddButton("Help", nil)
	
	// Quit button
	form.AddButton("Quit", func() {
		app.Stop()
	})

	form.AddTextView("Results:", result.switch_name, 0, 2, false, false)
	form.AddTextView("", "Bye", 0, 2, false, false)
	form.AddTextView("", "Hello", 0, 2, false, false)
	form.AddTextView("", "Goodbye", 0, 2, false, false)
	form.AddTextView("", "Greetings", 0, 2, false, false)

	err := app.SetRoot(form, true).Run();
	if err != nil {
		panic(err)
	}


	// Reset button
}

// Void
func save_link_data(data *link_data) {
	if data == nil {
		// Tui message pottentially
		log.Fatal("No link data is put")
	}
	curr_time := time.DateTime

	file, err := os.Create("LinkData_" + curr_time +".txt")
	if err != nil {
		log.Fatal("Error creating a file")
	}
	defer file.Close()
	file.WriteString("Switch name: " + (*data).switch_name +  "\n" + 
	"Port Identifier: " + (*data).port_id +  "\n" +
	"Vlan Identifier: " + (*data).vlan_id +  "\n" +
	"Switch Ip Address: " + (*data).switch_ip +  "\n" +
	"Switch model: " + (*data).switch_model +  "\n" +
	"Port Identifier: " + (*data).duplex_option +  "\n" +
	"Port Identifier: " + (*data).vpt_mgmt_domain +  "\n")

}

func get_link_data(connection *connect_struct) (link_data, bool) {

	var link link_data
	var result_found = false

	var device_argument string

	device_argument = (*(*connection).iface).Name

	device_argument = (*(*connection).iface).Name

	home := os.Getenv("HOME")

	cmd := exec.Command(
		"pkexec",
		home + "/personal/ldlinux/helper/helper",
		device_argument,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Pipe failed")
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal("Command not executed")
	}


	scanner := bufio.NewScanner(stdout)

	var line string 
	for scanner.Scan() {
		line = scanner.Text()

		var result LinkDataHelperForPackets

		err := json.Unmarshal([]byte(line), &result)
	
		if err != nil {
			continue
		}

		link.connection = connection
		link.switch_name = result.SwitchName
		link.switch_model = result.SwitchModel
		link.port_id = result.PortID
		link.switch_ip = result.SwitchIP
		link.vlan_id = result.VlanID
		link.vpt_mgmt_domain = result.VptMgmtDomain
		link.duplex_option = result.DuplexOption
		link.protocol = result.Protocol
		result_found = true

	}


	return link, result_found
}

// Excludes wireless for Windows and Linux however I mean its not like wireless works for switches anyway
func exclude_wireless(iface net.Interface) bool {
	if strings.Contains(iface.Name, "wi-fi") ||
		strings.Contains(iface.Name, "wifi") ||
		strings.Contains(iface.Name, "wireless") ||
		strings.Contains(iface.Name, "wlan") ||
		strings.HasPrefix(iface.Name, "wl") ||
		strings.HasPrefix(iface.Name, "wlan") {
		return false
	}

	return true
}

func ethernet_checker(ifaces []net.Interface) []valid_iface {
	var valid []valid_iface
	for _, v := range ifaces {
		var mock_valid valid_iface

		if v.Flags&net.FlagUp > 0 &&
			len(v.HardwareAddr) > 0 &&
			v.Flags&net.FlagLoopback == 0 &&
			exclude_wireless(v) {
			mock_valid.iface = &v
			if v.Flags&net.FlagRunning > 0 {
				mock_valid.is_running = true
				valid = append(valid, mock_valid)
			} else {
				valid = append(valid, mock_valid)

			}
		}
	}

	return valid
}

func get_connection_data() []connect_struct {
	// Get the Network card - implementation needed

	var connections []connect_struct
	var info connect_struct = connect_struct{
		2,
		"Ethernet",
		"",
		"",
		false,
		"",
		nil,
	}
	var info_false connect_struct = connect_struct{
		-1,
		"No Adapter",
		"<nil>",
		"<nil>",
		false,
		"<nil>",
		nil,
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to get network interface information")
		connections = append(connections, info_false)
	}

	valid_ifaces := ethernet_checker(ifaces)

	count := 0
	for _, v := range valid_ifaces {
		// https://www.golinuxcloud.com/golang-get-ip-address/
		addrs, err := v.iface.Addrs()
		if err != nil {
			fmt.Println(err)
		}

		var ip net.IP
		for _, addr := range addrs {
			// check the type of the address and
			// assign it to the variable ip of type net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
			case *net.IPNet:
				ip = v.IP
			default:
				continue
			}

			// IPv4
			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}

			if ip4.IsLoopback() || ip4.IsLinkLocalUnicast() {
				continue
			}

			ip = ip4
			break

		}
		// if strings.HasPrefix(ip.String(), "192.168") {
		// 	ipv4_str
		// }

		info.id = count + 1
		info.name = fmt.Sprintf("Ethernet %d", count)
		info.mac_addr = v.iface.HardwareAddr.String()
		info.network_card = v.iface.Name
		info.ip_addr = ip.String()

		if v.is_running && info.ip_addr != "<nil>" {
			info.status = true
			// Return the structure
		}
		info.iface = v.iface

		connections = append(connections, info)
		count++
	}
	return connections

}
