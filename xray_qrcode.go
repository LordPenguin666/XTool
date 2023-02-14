package main

import "github.com/Baozisoftware/qrcode-terminal-go"

func (c *Config) XrayLinkQRCode() *Config {
	obj := qrcodeTerminal.New().Get(c.XrayLink)
	obj.Print()
	return c
}
