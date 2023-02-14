package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"os/exec"
)

func (c *Config) CaddyInstallation() *Config {
	if c.CaddyVer != "" {
		fmt.Printf("%v %v %v",
			red("[Warning]"),
			yellow("Caddy 已存在, 是否重新安装?"),
			blue("(y|n)"),
		)

		for {
			var confirm string
			if _, err := fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				c.StopCaddy().InstallDefaultCaddy()
				break
			} else if confirm == "n" {
				break
			} else {
				fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
			}
		}

	}

	if c.CaddyProxyProtocolSupport {
		fmt.Printf("%v %v %v",
			red("[Warning]"),
			yellow("Caddy 已支持 Proxy Protocol, 是否重新安装?"),
			blue("(y|n)"),
		)

		for {
			var confirm string
			if _, err := fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				c.StopCaddy().ReplaceCaddyWithModules()
				break
			} else if confirm == "n" {
				break
			} else {
				fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
			}
		}
		return c
	}

	fmt.Printf("%v%v\n", green("预备 "), blue("执行全新安装"))

	// 安装 Caddy
	c.InstallDefaultCaddy().StopCaddy().ReplaceCaddyWithModules()

	// 配置 Caddy
	c.DeployCaddyFile().DeployCaddyConf().RestartCaddy().EnableCaddy()

	return c
}

func (c *Config) InstallDefaultCaddy() *Config {
	fmt.Printf("%v %v\n", green("系统发行版"), blue(c.Platform))

	switch c.Platform {
	case "debian":
		c.DebianCaddyInstall()
	case "rhel":
		switch c.PackageManagement {
		case "yum":
			c.ReadHatYumCaddyInstall()
		case "dnf":
			c.ReadHatDnfCaddyInstall()
		}

	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("系统不在支持的范围内!\n\n"))
		os.Exit(0)
	}

	return c
}

func (c *Config) DebianCaddyInstall() *Config {
	commands := []string{
		"sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https",
		"curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg",
		"curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list",
		"sudo apt update -y",
		"sudo apt -y install caddy -y",
	}

	for _, command := range commands {
		bash := exec.Command("bash", "-c", command)
		if output, err := bash.CombinedOutput(); err != nil {
			c.logger.Error(err.Error())
		} else {
			fmt.Println(string(output))
		}

	}

	return c
}

func (c *Config) ReadHatYumCaddyInstall() *Config {
	commands := []string{
		"yum install yum-plugin-copr -y",
		"yum copr enable @caddy/caddy -y",
		"yum install caddy -y",
	}

	for _, command := range commands {
		bash := exec.Command("bash", "-c", command)
		if _, err := bash.CombinedOutput(); err != nil {
			c.logger.Error(err.Error())
		}
	}

	return c
}

func (c *Config) ReadHatDnfCaddyInstall() *Config {
	commands := []string{
		"dnf install 'dnf-command(copr) -y",
		"dnf copr enable @caddy/caddy -y",
		"dnf install caddy -y",
	}

	for _, command := range commands {
		bash := exec.Command("bash", "-c", command)
		if _, err := bash.CombinedOutput(); err != nil {
			c.logger.Error(err.Error())
		}
	}

	return c
}

func (c *Config) ReplaceCaddyWithModules() *Config {
	params := map[string]string{
		"os": "linux",
		"p":  "github.com/mastercactapus/caddy2-proxyprotocol",
	}

	switch c.Arch {
	case "x86_64":
		params["arch"] = "amd64"
	//case "aarch64":
	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("系统不在支持的范围内!\n\n"))
		os.Exit(0)
	}

	client := resty.New()
	client.SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	_, err := client.R().
		SetQueryParams(params).
		SetOutput("/usr/bin/caddy").
		Get("https://caddyserver.com/api/download")

	if err != nil {
		fmt.Printf("%v %v", red("[Warning]"), yellow("Caddy 下载失败!\n\n"))
		c.logger.Fatal(err.Error())
	}

	return c
}
