package tui

import "github.com/rivo/tview"

func (ui *UI) menu() *tview.Grid {
	list := tview.NewList()

	list.AddItem("安装 Xray + Vless + xtls-rprx-vision + tls", "Xray 监听 443 并回落至 Caddy", '0', func() {
		ui.xrayFunc = ui.xrayVisionTLS
		ui.nextContent(ui.tlsForm())
	})

	list.AddItem("安装 Xray + Vless + tls-rprx-vision + reality", "免域名, 但是不够灵活 :(", '1', func() {

	})

	list.AddItem("Caddy 配置", "", '2', func() {
		ui.nextContent(ui.caddyMenu())
	})

	list.AddItem("Xray 配置", "", '3', func() {
		ui.nextContent(ui.xrayMenu())
	})

	list.AddItem("Xray QR Code", "生成 Xray 二维码", '4', func() {
		ui.nextContent(ui.xrayQrcode())
	})

	list.AddItem("Exit", "退出程序", 'q', func() {
		ui.app.Stop()
	})

	return ui.contentWithTitle("Menu", list)
}
