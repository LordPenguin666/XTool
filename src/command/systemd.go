package command

import "os/exec"

func SystemdReload() *exec.Cmd {
	return exec.Command("systemctl", "daemon-reload")
}
