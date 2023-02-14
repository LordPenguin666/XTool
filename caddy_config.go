package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (c *Config) DeployCaddyFile() *Config {
	caddyFile := "{\n" +
		fmt.Sprintf("        servers :%v {\n", c.CaddyHTTPPort) +
		"                listener_wrappers {\n" +
		"                        proxy_protocol {\n" +
		"                                timeout 2s\n" +
		"                                allow 0.0.0.0/0\n" +
		"                        }\n" +
		"                        tls\n" +
		"                }\n" +
		"                protocols h1 h2 h2c h3\n" +
		"        }\n" +
		"}\n" +
		"\n" +
		":80 {\n" +
		"    redir https://{host}{url}\n" +
		"}\n" +
		"\n" +
		"import /etc/caddy/conf.d/*"

	command := exec.Command("cat", "/etc/caddy/Caddyfile")
	catCaddyFile, err := command.CombinedOutput()
	if err != nil {
		c.logger.Error(err.Error())
	}

	if c.FileExist("/etc/caddy/Caddyfile") && strings.Contains(string(catCaddyFile), "proxy_protocol") {
		fmt.Printf("%v\n", green("当前的 Caddy 配置:"))
		fmt.Printf("%v\n", string(catCaddyFile))
		fmt.Printf("%v %v %v",
			red("[Warning]"),
			yellow("Caddyfile 可能已经配置完毕, 是否覆盖?"),
			blue("(y|n)"),
		)

		for {
			var confirm string
			if _, err = fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				break
			} else if confirm == "n" {
				return c
			} else {
				fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
			}
		}
	}

	if err = os.WriteFile("/etc/caddy/Caddyfile", []byte(caddyFile), 644); err != nil {
		c.logger.Error(err.Error())
	}

	return c
}

func (c *Config) DeployCaddyConf() *Config {
	createDir := exec.Command("mkdir", "-p", "/etc/caddy/conf.d/")
	if err := createDir.Run(); err != nil {
		c.logger.Error(err.Error())
	}

	confPath := "/etc/caddy/conf.d/" + c.Domain
	if c.FileExist(confPath) {
		command := exec.Command("cat", confPath)
		catCaddyConf, err := command.CombinedOutput()
		if err != nil {
			c.logger.Error(err.Error())
		}

		fmt.Printf("%v\n", green("当前的 Caddy Conf 配置:"))
		fmt.Printf("%v\n", string(catCaddyConf))
		fmt.Printf("%v %v %v",
			red("[Warning]"),
			yellow("Caddy Conf 可能已经配置完毕, 是否覆盖?"),
			blue("(y|n)"),
		)

		for {
			var confirm string
			if _, err = fmt.Scan(&confirm); err != nil {
				c.logger.Error(err.Error())
			}
			if confirm == "y" {
				break
			} else if confirm == "n" {
				return c
			} else {
				fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
			}
		}
	}

	var caddyConf string

	if c.isHttps {
		caddyConf = fmt.Sprintf("http://%v:%v {\n", c.Domain, c.CaddyHTTPPort) +
			fmt.Sprintf("    reverse_proxy %v {\n", c.ProxyDomain) +
			"        header_up Host {upstream_hostport}\n" +
			"        transport http {\n" +
			"            tls\n" +
			"        }\n" +
			"    }\n" +
			"}\n" +
			"\n" +
			fmt.Sprintf("%v:%v {\n", c.Domain, c.CaddyHTTPSPort) +
			fmt.Sprintf("    reverse_proxy %v {\n", c.ProxyDomain) +
			"        header_up Host {upstream_hostport}\n" +
			"        transport http {\n" +
			"            tls\n" +
			"        }\n" +
			"    }\n" +
			"}"
	} else {
		caddyConf = fmt.Sprintf("http://%v:%v {\n", c.Domain, c.CaddyHTTPPort) +
			fmt.Sprintf("    reverse_proxy %vn", c.ProxyDomain) +
			"}\n" +
			"\n" +
			fmt.Sprintf("%v:%v {\n", c.Domain, c.CaddyHTTPSPort) +
			fmt.Sprintf("    reverse_proxy %v\n", c.ProxyDomain) +
			"}"
	}

	if err := os.WriteFile(confPath, []byte(caddyConf), 644); err != nil {
		c.logger.Error(err.Error())
	}

	return c
}
