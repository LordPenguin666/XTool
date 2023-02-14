package main

import (
	"fmt"
	"github.com/fatih/color"
)

func (c *Config) menu() {
	color.Blue("********************************************************************")
	color.Blue("*                                                                  *")

	if _, err := color.New(color.FgBlue).Print("*"); err != nil {
		c.logger.Error(err.Error())
	}

	if _, err := color.New(color.FgRed).
		Print("                         Xray tool v1.0                           "); err != nil {
		c.logger.Error(err.Error())
	}

	if _, err := color.New(color.FgBlue).Println("*"); err != nil {
		c.logger.Error(err.Error())
	}

	color.Blue("*             https://github.com/LordPenguin666/XTool              *")

	if _, err := color.New(color.FgBlue).Print("*"); err != nil {
		c.logger.Error(err.Error())
	}

	if _, err := color.New(color.FgYellow).
		Print("          给面包狗点个 Star 喵, 给面包狗点个 Star 谢谢喵          "); err != nil {
		c.logger.Error(err.Error())
	}

	if _, err := color.New(color.FgBlue).Println("*"); err != nil {
		c.logger.Error(err.Error())
	}

	color.Blue("*                                                                  *")
	color.Blue("********************************************************************")

	fmt.Println()

	// Xray version (key)
	if _, err := color.New(color.FgGreen).Print("\tXray: "); err != nil {
		c.logger.Error(err.Error())
	}

	// xray version (value)
	if c.XrayVer == "" {
		if _, err := color.New(color.FgRed).Print("未安装"); err != nil {
			c.logger.Error(err.Error())
		}
	} else {
		if _, err := color.New(color.FgBlue).Print("v" + c.XrayVer); err != nil {
			c.logger.Error(err.Error())
		}
	}

	// Caddy version (key)
	if _, err := color.New(color.FgGreen).Print("\tCaddy: "); err != nil {
		c.logger.Error(err.Error())
	}

	// Caddy version (value)
	if c.CaddyVer != "" {
		if _, err := color.New(color.FgBlue).Print(c.CaddyVer); err != nil {
			c.logger.Error(err.Error())
		}
	} else {
		if _, err := color.New(color.FgRed).Print("未安装"); err != nil {
			c.logger.Error(err.Error())
		}
	}

	fmt.Print("\t")

	// Caddy proxy protocol support (key)
	if _, err := color.New(color.FgGreen).Print("Proxy Protocol: "); err != nil {
		c.logger.Error(err.Error())
	}

	// Caddy proxy protocol support (value)
	if c.CaddyProxyProtocolSupport {
		if _, err := color.New(color.FgBlue).Print("支持"); err != nil {
			c.logger.Error(err.Error())
		}
	} else {
		if _, err := color.New(color.FgRed).Print("不支持"); err != nil {
			c.logger.Error(err.Error())
		}
	}

	fmt.Println()
	fmt.Println()

	fmt.Printf("  %v %v\n", green("1) "), blue("安装 xray + vless + xtls-rprx-vision + Caddy 回落 (默认)"))
	fmt.Printf("  %v %v\n", green("2) "), blue("安装 xray + vless + xtls-rprx-vision + navieproxy 回落 + Caddy 回落 (在写了在写了)"))
	fmt.Println()
	fmt.Printf("  %v %v\n", green("11) "), blue("修改 Caddy HTTP 端口 (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("12) "), blue("修改 Caddy HTTPS 端口 (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("13) "), blue("启动 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("14) "), blue("重启 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("15) "), blue("关闭 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("16) "), blue("自启动 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("17) "), blue("关闭自启动 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("18) "), blue("卸载 Caddy (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("19) "), blue("安装 Caddy (在写了在写了)"))
	fmt.Println()
	fmt.Printf("  %v %v\n", green("21) "), blue("修改 Xray Xtls 端口 (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("22) "), blue("修改 Xray naiveproxy 端口 (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("23) "), blue("启动 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("24) "), blue("重启 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("25) "), blue("关闭 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("26) "), blue("自启动 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("27) "), blue("关闭自启动 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("28) "), blue("卸载 Xray (在写了在写了)"))
	fmt.Printf("  %v %v\n", green("29) "), blue("安装 Xray (在写了在写了)"))
	fmt.Println()
	fmt.Printf("  %v %v\n", green("0) "), blue("退出此脚本"))
	fmt.Println()
}
