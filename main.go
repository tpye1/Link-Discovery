package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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
	//stacked	bool
	//member	int
}

func main() {

	if runtime.GOOS == "windows" {
		if !win_isAdmin() {
			log.Fatal("Please run this application from an Administrator PowerShell/CMD window.")
		}

	}

	connections := get_connection_data()

	p := tea.NewProgram(
		initialModel(connections),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// Void
func save_link_data(data *link_data) error {
	var err error = nil
	if data == nil {
		// Tui message pottentially
		err = errors.New("No link data is put")
		return err
	}
	curr_time := time.DateTime

	file, err := os.Create("LinkData_" + curr_time + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("Switch name: " + (*data).switch_name + "\n" +
		"Port Identifier: " + (*data).port_id + "\n" +
		"Vlan Identifier: " + (*data).vlan_id + "\n" +
		"Switch Ip Address: " + (*data).switch_ip + "\n" +
		"Switch model: " + (*data).switch_model + "\n" +
		"Port Identifier: " + (*data).duplex_option + "\n" +
		"Port Identifier: " + (*data).vpt_mgmt_domain + "\n")
	return err

}

func get_link_data(connection *connect_struct) (link_data, error) {

	var link link_data

	var device_argument string
	if runtime.GOOS == "windows" {
		device, err := windows_pcap_translate(connection.iface.Name)
		if err != nil {
			return link, err
		}
		var handle *pcap.Handle

		handle, err = pcap.OpenLive(device, 1600, true, 1)
		if err != nil {
			return link, err
		}
		defer handle.Close()

		err = handle.SetBPFFilter(
			"(ether proto 0x88cc) or " +
				"(ether[12:2] <= 1500 and ether[14:1] == 0xaa and ether[15:1] == 0xaa and ether[16:1] == 0x03 and ether[20:2] == 0x2000)",
		)

		if err != nil {
			return link, err
		}

		var packet_source *gopacket.PacketSource
		timeout := time.After(time.Minute)

		packet_source = gopacket.NewPacketSource(handle, handle.LinkType())
		for {
			select {
			case packet := <-packet_source.Packets():
				if cdp := packet.Layer(layers.LayerTypeCiscoDiscovery); cdp != nil {
					link = get_cdp(cdp, connection)
				}
				if lldp := packet.Layer(layers.LayerTypeLinkLayerDiscovery); lldp != nil {
					link = get_lldp(lldp, connection)
				}

			case <-timeout:
				err = errors.New("Timeout")
				return link, err
			}

		}

	} else {

		device_argument = (*(*connection).iface).Name

		home := os.Getenv("HOME")

		cmd := exec.Command(
			"pkexec",
			home+"/ldlinux/helper/helper",
			device_argument,
		)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return link, err
		}
		err = cmd.Start()
		if err != nil {
			return link, err
		}

		scanner := bufio.NewScanner(stdout)

		var line string
		for scanner.Scan() {
			line = scanner.Text()

			var result LinkDataHelperForPackets

			err := json.Unmarshal([]byte(line), &result)

			if err != nil {
				fmt.Errorf("JSON ERROR: %v\n", err)
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
			break
		}
		cmd.Process.Kill()
	}
	return link, nil
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

	for i := range ifaces {
		iface := &ifaces[i]

		if iface.Flags&net.FlagUp > 0 &&
			len(iface.HardwareAddr) > 0 &&
			iface.Flags&net.FlagLoopback == 0 &&
			exclude_wireless(*iface) {
			mockValid := valid_iface{
				iface: iface,
			}

			if iface.Flags&net.FlagRunning > 0 {
				mockValid.is_running = true
			}

			valid = append(valid, mockValid)
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
