package tui

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
)

func (ui *UI) xrayQrcode() *tview.Grid {
	xrayRawConfig, err := os.ReadFile("/usr/local/etc/xray/config.json")
	if err != nil {
		panic(err)
	}

	u := url.Values{}
	security := gjson.Get(string(xrayRawConfig), "inbounds").Array()[0].Get("streamSettings").Get("security").String()

	switch security {
	case "tls":
		xconfig := &XrayTLSConfig{}
		if err = json.Unmarshal(xrayRawConfig, xconfig); err != nil {
			panic(err)
		}

		localDomain, ok := ui.getLocalDomain()
		if !ok {
			ui.nextContent(ui.menu())
		}

		u.Add("security", "tls")
		u.Add("encryption", "none")
		u.Add("alpn", "h2,http/1.1")
		u.Add("headerType", "none")
		u.Add("fp", "random")
		u.Add("type", "tcp")
		u.Add("flow", "xtls-rprx-vision")
		u.Add("sni", localDomain.Hostname())

		name := fmt.Sprintf("[%v] (%v) xtls-rprx-vision", localDomain.Hostname(), ui.addr)

		XrayLink := fmt.Sprintf("vless://%v@%v:%v?%v#%v",
			xconfig.Inbounds[0].Settings.Clients[0].Id, localDomain.Hostname(), 443, u.Encode(), name)

		qr, err := qrcode.New(XrayLink, qrcode.Medium)
		if err != nil {
			panic(err)
		}

		textView := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true).
			SetText("\n" + qr.ToSmallString(false))

		return ui.contentWithTitle("Xray QR Code", textView)

	default:
		panic("unsupported security")
	}
}
