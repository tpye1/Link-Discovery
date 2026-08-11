//go:build windows

package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/windows"
)

func win_isAdmin() bool {
	token := windows.GetCurrentProcessToken()

	elevated := token.IsElevated()
	return elevated

}

func windows_pcap_translate(iface_name string) (string, error) {

	// Find all devices
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}

	for _, device := range devs {
		if strings.Contains(strings.ToLower(device.Description),
			strings.ToLower(iface_name),
		) {
			return device.Name, nil
		}

	}
	return "", fmt.Errorf("No pcap device found!")
}

func get_lldp(lldp_layer gopacket.Layer, connection *connect_struct) link_data {

	var lldp_link_data = link_data{
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
				continue
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

	var cdp_link_data = link_data{
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
				continue
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
