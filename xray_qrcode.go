package main

import (
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
)

func (c *Config) XrayLinkQRCode() *Config {
Loop:
	fmt.Printf("%v%v%v%v",
		green("是否生成 vless 二维码"),
		blue("(y|n): "))

	var input string
	if _, err := fmt.Scan(&input); err != nil {
		c.logger.Error(err.Error())
	}

	switch input {
	case "y":
		break
	case "n":
		return c
	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的选项!\n\n"))
		goto Loop
	}

	fmt.Println()
	fmt.Println()
	obj := qrcodeTerminal.New().Get(c.XrayLink)
	obj.Print()

	fmt.Println()
	return c
}
