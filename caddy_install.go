package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"os"
	"os/exec"
	"time"
)

func (c *Config) CaddyInstallation(conf bool) *Config {
	var skipInstallCaddy, skipReplaceCaddy bool

	if c.CaddyVer != "" {
		fmt.Printf("%v %v %v",
			red("[Warning]"),
			yellow("Caddy 已存在, 是否重新安装?"),
			blue("(y|n): "),
		)

		for {
			var confirm string
			if _, err := fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				c.StopCaddy()
				break
			} else if confirm == "n" {
				skipInstallCaddy = true
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
			blue("(y|n): "),
		)

		for {
			var confirm string
			if _, err := fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				c.StopCaddy()
				skipReplaceCaddy = true
				break
			} else if confirm == "n" {
				break
			} else {
				fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
			}
		}
		return c
	}

	fmt.Printf("%v %v\n", green("系统发行版"), blue(c.Platform))
	fmt.Printf("%v%v\n", green("预备执行安装"))

	// 安装 Caddy
	c.InstallDefaultCaddy(skipInstallCaddy).StopCaddy().ReplaceCaddyWithModules(skipReplaceCaddy)

	if !conf {
		return c
	}

	// 配置 Caddy
	c.DeployCaddyFile().DeployCaddyConf().RestartCaddy().EnableCaddy()

	// 等待证书生成
	fmt.Printf("%v\n", "等待 Caddy 申请证书中...")
	now := time.Now()
	timeLeft := 30 * time.Second

	for {
		task := time.NewTimer(50 * time.Millisecond)
		if time.Since(now) >= timeLeft {
			break
		}

		t := timeLeft - time.Since(now)
		fmt.Printf("%v %v %v\r",
			green("剩余时间:"),
			blue(decimal.NewFromFloat(t.Seconds()).Round(1).String()),
			green("s"),
		)

		<-task.C
	}

	// 检查证书
	c.CaddySSL()

	return c
}

func (c *Config) InstallDefaultCaddy(skip bool) *Config {
	if skip {
		return c
	}

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
		fmt.Printf("%v %v\n", green("执行命令"), blue(command))
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

func (c *Config) ReplaceCaddyWithModules(skip bool) *Config {
	if skip {
		return c
	}

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

func (c *Config) UninstallCaddy() *Config {
	switch c.PackageManagement {
	case "apt":
		bash := exec.Command("apt", "autoremove", "caddy", "-y")
		output, err := bash.CombinedOutput()
		if err != nil {
			c.logger.Error(err.Error())
		}
		fmt.Println(output)
	case "yum":
		bash := exec.Command("yum", "remove", "caddy", "-y")
		output, err := bash.CombinedOutput()
		if err != nil {
			c.logger.Error(err.Error())
		}
		fmt.Println(output)
	case "dnf":
		bash := exec.Command("dnf", "remove", "caddy", "-y")
		output, err := bash.CombinedOutput()
		if err != nil {
			c.logger.Error(err.Error())
		}
		fmt.Println(output)
	}
	return c
}
