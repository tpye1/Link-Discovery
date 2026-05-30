package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Link Discovery Client for Linux",
		Width:  850,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func get_link_data(connection *connect_struct) link_data {

	var link link_data

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

	}


	return link
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
	// Get the Network card - implenentation needed

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
