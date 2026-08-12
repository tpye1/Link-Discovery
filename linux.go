//go:build !windows

package main

import (
	"fmt"

	"github.com/google/gopacket"
)

// This is to make the compiler not scream at me work even though it is unreachable code. Ai helped with this icl.

func win_isAdmin() bool {
	return true
}

func relaunchAsAdmin() error {
	return nil
}

func windows_pcap_translate() (string, error) {

	return "", nil
}

func get_lldp(lldp_layer gopacket.Layer, connection *connect_struct) link_data {
	fmt.Println(lldp_layer)
	var link link_data
	link.connection = connection
	return link
}

func get_cdp(cdp_layer gopacket.Layer, connection *connect_struct) link_data {
	fmt.Println(cdp_layer)
	var link link_data
	link.connection = connection
	return link
}
