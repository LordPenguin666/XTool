package main

import "fmt"

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
