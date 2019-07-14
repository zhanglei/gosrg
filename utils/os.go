package utils

import (
	"os/exec"
	"runtime"
)

func OpenLink(link string) {
	command := "open"
	if runtime.GOOS == "linux" {
		command = "x-www-browser"
	}
	if _, err := RunCommand(command, link); err != nil {
		Logger.Printf("OpenLink error: %s\n", err)
	}
}

func RunCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	res, err := cmd.Output()
	if err != nil {
		Logger.Printf("command error: %s\n", err)
		return []byte{}, nil
	}
	return res, nil
}
