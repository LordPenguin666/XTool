package main

import (
	"fmt"
	"strings"
)

func (c *Config) ConfirmModify() *Config {
Loop:
	var input string

	fmt.Printf("%v%v%v%v",
		green("是否修改默认配置"),
		blue(" (端口号)"),
		yellow(" [不建议修改] "),
		red("(y|n): "))
	if _, err := fmt.Scan(&input); err != nil {
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
		goto Loop
	}

	switch input {
	case "y":
		c.ModifyConfig()
	case "n":
		return c
	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
		goto Loop
	}

	c.AddDomain().DomainLookup().AddProxyDomain()

	return c
}

func (c *Config) ModifyConfig() *Config {
	var (
		httpPort  int
		httpsPort int
	)

	for {
		fmt.Printf("%v%v%v: ",
			green("请输入 Caddy 监听的"),
			blue(" [HTTP 端口]"),
			yellow(" (默认 8080)"),
		)
		if _, err := fmt.Scan(&httpPort); err != nil {
			fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的数字!\n\n"))
			continue
		}
		c.CaddyHTTPPort = httpPort
		break
	}

	for {
		fmt.Printf("%v%v%v: ",
			green("请输入 Caddy 监听的"),
			blue(" [HTTPS 端口] "),
			yellow("(默认 8443)"),
		)
		if _, err := fmt.Scan(&httpsPort); err != nil {
			fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的数字!\n\n"))
			continue
		}
		c.CaddyHTTPSPort = httpsPort
		break
	}

	return c
}

func (c *Config) AddDomain() *Config {
	var input string

	fmt.Printf("%v %v",
		green("请输入解析到此服务器的域名, 无需 http 前缀"),
		blue(" (例如: xxx.example.com)"),
	)

	if _, err := fmt.Scan(&input); err != nil {
		c.logger.Error(err.Error())
	}

	c.Domain = input
	return c
}

func (c *Config) AddProxyDomain() *Config {
	var input string

Loop:
	fmt.Printf("%v %v",
		green("请输入要反代的域名, 需要 http 前缀"),
		blue(" (例如: https://bing.com)"),
	)

	if _, err := fmt.Scan(&input); err != nil {
		c.logger.Error(err.Error())
	}

	c.ProxyDomain = input

	words := strings.Split(c.ProxyDomain, ":")

	switch words[0] {
	case "http":
		c.isHttps = false
	case "https":
		c.isHttps = true
	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的反代域名!\n\n"))
		goto Loop
	}

	return c
}
