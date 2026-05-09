package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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

var port_duplex = [3]string{"Full", "Half", "<nil>"}

type valid_iface struct {
	iface      net.Interface
	is_running bool
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

// Require sudo privilages

func get_lldp(lldp_layer gopacket.Layer, connection *connect_struct) link_data {

	var lldp_link_data link_data = link_data{
		connection,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"LLDP",
	}

	lldp := lldp_layer.(*layers.LinkLayerDiscovery)
	for _, v := range lldp.Values {
		switch v.Type {
		case layers.LLDPTLVSysName:
			lldp_link_data.switch_name = string(v.Value)
		case layers.LLDPTLVSysDescription:
			lldp_link_data.switch_model = string(v.Value)
		case layers.LLDPTLVPortDescription:
			lldp_link_data.port_id = string(v.Value)
		case layers.LLDPTLVMgmtAddress:
			if len(v.Value) < 6 {
				log.Fatal("Stop")
			}
			addr_subtype := v.Value[1]
			if addr_subtype == 1 {
				ipv4 := net.IPv4(v.Value[2], v.Value[3], v.Value[4], v.Value[5])
				lldp_link_data.switch_ip = ipv4.String()
			}
		}
	}

	return lldp_link_data
}

func get_cdp(cdp_layer gopacket.Layer, connection *connect_struct) link_data {

	var cdp_link_data link_data = link_data{
		connection,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"CDP",
	}

	cdp := cdp_layer.(*layers.CiscoDiscovery)
	for _, v := range cdp.Values {
		switch v.Type {
		case layers.CDPTLVDevID:
			cdp_link_data.switch_name = string(v.Value)
		case layers.CDPTLVPortID:
			cdp_link_data.port_id = string(v.Value)
		case layers.CDPTLVNativeVLAN:
			cdp_link_data.vlan_id = string(v.Value)
		case layers.CDPTLVAddress:
			if len(v.Value) < 6 {
				log.Fatal("Stop")
			}
			addr_subtype := v.Value[1]
			if addr_subtype == 1 {
				ipv4 := net.IPv4(v.Value[2], v.Value[3], v.Value[4], v.Value[5])
				cdp_link_data.switch_ip = ipv4.String()
			}
		case layers.CDPTLVPlatform:
			cdp_link_data.switch_name = string(v.Value)
		case layers.CDPTLVFullDuplex:
			cdp_link_data.duplex_option = string(v.Value)
		case layers.CDPTLVVTPDomain:
			cdp_link_data.vpt_mgmt_domain = string(v.Value)
		}
	}
	return cdp_link_data
}

func windows_pcap_translate(iface *net.Interface) (string, error) {

	// Find all devices
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}

	for _, device := range devs {
		if strings.Contains(strings.ToLower(device.Description),
			strings.ToLower(iface.Name),
		) {
			return device.Name, nil
		}

	}
	return "", fmt.Errorf("No pcap device found!")
}

func get_link_data(connection *connect_struct) link_data {
	var err error
	var handle *pcap.Handle
	// Ethernet Connection,
	// Number of bytes taken from each packet 1600 is a sufficient number,
	// promiscous will be enabled,
	// We will choose 1 milisecond as soon as a packet is found it will stop

	var device string
	if runtime.GOOS == "windows" {
		device, err = windows_pcap_translate((*connection).iface)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		device = connection.network_card
	}

	handle, err = pcap.OpenLive(device, 1600, true, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.SetBPFFilter(
		"(ether proto 0x88cc) or " +
			"(ether[12:2] <= 1500 and ether[14:1] == 0xaa and ether[15:1] == 0xaa and ether[16:1] == 0x03 and ether[20:2] == 0x2000)",
	)

	if err != nil {
		log.Fatal(err)
	}

	var packet_source *gopacket.PacketSource
	timeout := time.After(time.Minute)

	packet_source = gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case packet := <-packet_source.Packets():
			if cdp := packet.Layer(layers.LayerTypeCiscoDiscovery); cdp != nil {
				return get_cdp(cdp, connection)
			}
			if lldp := packet.Layer(layers.LayerTypeLinkLayerDiscovery); lldp != nil {
				return get_lldp(lldp, connection)
			}

		case <-timeout:
			log.Fatal("Waited too long")
		}

	}
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
			mock_valid.iface = v
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
		info.iface = &v.iface

		connections = append(connections, info)
		count++
	}
	return connections

}
