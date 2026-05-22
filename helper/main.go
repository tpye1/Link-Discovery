package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No Arguments please add your interface")
	}

	device := os.Args[1]
	admin_helper(device)
}


func admin_helper(device string)  {

	var err error
	var handle *pcap.Handle
	// Ethernet Connection,
	// Number of bytes taken from each packet 1600 is a sufficient number,
	// promiscous will be enabled,
	// We will choose 1 milisecond as soon as a packet is found it will stop

	var link LinkDataHelperForPackets 

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
				link = get_cdp(cdp)
				jsonData, err := json.Marshal(link)
				if err != nil {
					continue
				}
				fmt.Println(string(jsonData))
			}
			if lldp := packet.Layer(layers.LayerTypeLinkLayerDiscovery); lldp != nil {
				link = get_lldp(lldp)
				jsonData, err := json.Marshal(link)
				if err != nil {
					continue
				}
				fmt.Println(string(jsonData))
			}

		case <-timeout:
			log.Fatal("Waited too long")
		}

	}
}

func get_lldp(lldp_layer gopacket.Layer) LinkDataHelperForPackets {

	var lldp_link_data LinkDataHelperForPackets

	lldp_link_data.Protocol = "LLDP"

	lldp := lldp_layer.(*layers.LinkLayerDiscovery)
	for _, v := range lldp.Values {
		switch v.Type {
		case layers.LLDPTLVSysName:
			lldp_link_data.SwitchName = string(v.Value)
		case layers.LLDPTLVSysDescription:
			lldp_link_data.SwitchModel = string(v.Value)
		case layers.LLDPTLVPortID:
			lldp_link_data.PortID = string(v.Value)
		case layers.LLDPTLVMgmtAddress:
			if len(v.Value) < 6 {
				continue
			}
			addr_subtype := v.Value[1]
			if addr_subtype == 1 {
				ipv4 := net.IPv4(v.Value[2], v.Value[3], v.Value[4], v.Value[5])
				lldp_link_data.SwitchIP = ipv4.String()
			}
		}
	}
	return lldp_link_data
}

func get_cdp(cdp_layer gopacket.Layer) LinkDataHelperForPackets {

	var cdp_link_data LinkDataHelperForPackets
	cdp_link_data.Protocol = "CDP"


	cdp := cdp_layer.(*layers.CiscoDiscovery)
	for _, v := range cdp.Values {
		switch v.Type {
		case layers.CDPTLVDevID:
			cdp_link_data.SwitchName = string(v.Value)
		case layers.CDPTLVPortID:
			cdp_link_data.PortID = string(v.Value)
		case layers.CDPTLVNativeVLAN:
			cdp_link_data.VlanID = string(v.Value)
		case layers.CDPTLVAddress:
			if len(v.Value) < 6 {
				continue
			}
			addr_subtype := v.Value[1]
			if addr_subtype == 1 {
				ipv4 := net.IPv4(v.Value[2], v.Value[3], v.Value[4], v.Value[5])
				cdp_link_data.SwitchIP = ipv4.String()
			}
		case layers.CDPTLVPlatform:
			cdp_link_data.SwitchName = string(v.Value)
		case layers.CDPTLVFullDuplex:
			cdp_link_data.DuplexOption = string(v.Value)
		case layers.CDPTLVVTPDomain:
			cdp_link_data.VptMgmtDomain = string(v.Value)
		}
	}
	return cdp_link_data
}
