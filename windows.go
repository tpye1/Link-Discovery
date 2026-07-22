//go:build windows

package main

import (
	"os"
	"os/exec"

	"golang.org/x/sys/windows"
)

func win_isAdmin() bool {
	token := windows.Token(0)

	elevated := token.IsElevated()
	return elevated

}

func relaunchAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	return exec.Command(
		"powershell",
		"-Command",
		"Start-Process",
		exe,
		"-Verb",
		"RunAs",
	).Run()
}
