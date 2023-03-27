package command

import (
	"fmt"
	"os/exec"
)

func CaddyVersion() *exec.Cmd {
	return exec.Command("caddy", "version")
}

func CaddyListModules() *exec.Cmd {
	return exec.Command("caddy", "list-modules")
}

func CaddySystemd(command string) *exec.Cmd {
	return exec.Command("systemctl", command, "caddy")
}

func CaddyConfigMkdir() *exec.Cmd {
	return exec.Command("mkdir", "-p", "/etc/caddy/conf.d")
}

func CaddyRedHatInstall() []*exec.Cmd {
	return []*exec.Cmd{
		exec.Command("dnf", "install", "dnf-command(copr)", "-y"),
		exec.Command("dnf", "copr", "enable", "@caddy/caddy", "-y"),
		exec.Command("dnf", "install", "caddy ", "-y"),
	}
}

func CaddyDebianInstall() []*exec.Cmd {
	return []*exec.Cmd{
		exec.Command("bash", "-c", "sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https"),
		exec.Command("bash", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --batch --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg"),
		exec.Command("bash", "-c", "curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list"),
		exec.Command("bash", "-c", "sudo apt update -y"),
		exec.Command("bash", "-c", "sudo apt -y install caddy -y"),
	}
}

func CaddyRedHatUninstall() *exec.Cmd {
	return exec.Command("bash", "-c", "dnf autoremove caddy -y")
}

func CaddyDebianUninstall() *exec.Cmd {
	return exec.Command("bash", "-c", "apt autoremove caddy -y")
}

//func CaddyArchInstall() []*exec.Cmd {
//	return []*exec.Cmd{
//		exec.Command("Pacman", "-Sy", "caddy --noconfirm --needed"),
//	}
//}

func CaddyKeyPath(path, domain, suffix string) *exec.Cmd {
	return exec.Command("bash", "-c",
		fmt.Sprintf("find %v -name %v.%v", path, domain, suffix),
	)
}

//func CaddyUserAdd() *exec.Cmd {
//	return exec.Command("useradd", "caddy")
//}

func CaddyCertPath() string {
	return "/var/lib/caddy/certificates"
}

func CaddyRootCertPath() string {
	return "/root/.local/share/caddy/certificates"
}

func CaddyShareCertPath() string {
	return "/var/lib/caddy/.local/share/caddy/certificates"
}
