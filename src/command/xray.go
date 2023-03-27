package command

import "os/exec"

func XrayVersion() *exec.Cmd {
	return exec.Command("xray", "version")
}

func XraySystemd(command string) *exec.Cmd {
	return exec.Command("systemctl", command, "xray")
}

func XrayInstall() *exec.Cmd {
	return exec.Command("bash", "-c", "./install-release.sh install -u caddy")
}

func XrayUninstall() *exec.Cmd {
	return exec.Command("bash", "-c", "./install-release.sh remove")
}
