//go:build !windows

package main

// This is to make the compiler not scream at me work even though it is unreachable code. Ai helps with this icl.

func win_isAdmin() bool {
	return true
}

func relaunchAsAdmin() error {
	return nil
}

func windows_pcap_translate(iface_name string) (string, error) {

	return "", nil
}
