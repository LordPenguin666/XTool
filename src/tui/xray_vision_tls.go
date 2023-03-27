package tui

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"net"
	"net/url"
	"os"
	"strconv"
	"xray-tool/src/command"
)

func (ui *UI) tlsForm() *tview.Grid {
	form := tview.NewForm()

	form.AddInputField("Your Domain", ui.localDomain.Hostname(), 30, nil, nil)
	form.AddInputField("Proxy URI", "https://www.lovelive-anime.jp/", 50, nil, nil)
	form.AddInputField("Xray UUID", ui.xrayUUID, 50, nil, nil)
	form.AddInputField("Xray Port", "443", 10, nil, nil)
	form.AddInputField("Caddy Proxy Protocol Port", "8080", 10, nil, nil)
	form.AddInputField("Caddy HTTPS Port", "8443", 10, nil, nil)

	form.AddTextView("Example", "Your Domain:   us.ouo.eu.org\nProxy URI:     https://www.lovelive-anime.jp/", 50, 2, true, true)

	form.AddButton("Save", func() {
		var err error

		yourDomain := form.GetFormItemByLabel("Your Domain").(*tview.InputField).GetText()
		u, err := url.Parse(yourDomain)
		if u.Scheme == "http" {
			ui.localDomain = u
		} else {
			if ui.localDomain, err = url.Parse(fmt.Sprintf("https://%v", yourDomain)); err != nil {
				panic(err)
			}
		}

		if ui.proxyURI, err = url.Parse(form.GetFormItemByLabel("Proxy URI").(*tview.InputField).GetText()); err != nil {
			panic(err)
		}

		if ui.xrayPort, err = strconv.Atoi(form.GetFormItemByLabel("Xray Port").(*tview.InputField).GetText()); err != nil {
			panic(err)
		}

		if ui.caddyProxyProtocolPort, err = strconv.Atoi(form.GetFormItemByLabel("Caddy Proxy Protocol Port").(*tview.InputField).GetText()); err != nil {
			panic(err)
		}

		if ui.caddyHTTPSPort, err = strconv.Atoi(form.GetFormItemByLabel("Caddy HTTPS Port").(*tview.InputField).GetText()); err != nil {
			panic(err)
		}

		ui.xrayUUID = form.GetFormItemByLabel("Xray UUID").(*tview.InputField).GetText()

		ui.nextContent(ui.tlsConfirm())
	})

	form.AddButton("Cancel", func() {
		ui.nextContent(ui.menu())
	})

	return ui.contentWithTitle("Basic Info", form)
}

func (ui *UI) tlsConfirm() *tview.Grid {
	ip, err := net.LookupIP(ui.localDomain.Hostname())
	if err != nil {
		panic(err)
	}

	//ns, err := net.LookupNS(ui.localDomain.Hostname())
	//if err != nil {
	//	panic(err)
	//}

	var ipText string
	var correct bool

	for i := 0; i < len(ip); i++ {
		if ip[i].String() == ui.addr {
			correct = true
			ipText += fmt.Sprintf("[green]%v[white]", ip[i].String())
		} else {
			ipText += fmt.Sprintf("[red]%v[white]", ip[i].String())
		}

		if i != len(ip)-1 {
			ipText += "\n"
		}
	}

	//for i := 0; i < len(ns); i++ {
	//	nsText += ns[i].Host
	//	if i != len(ns)-1 {
	//		nsText += "\n"
	//	}
	//}

	var addr string
	if correct {
		addr = fmt.Sprintf("[green]%v[white]", ui.addr)
	} else {
		addr = fmt.Sprintf("[red]%v[white]", ui.addr)
	}

	form := tview.NewForm()
	form.AddTextView("Local Addr", addr, 50, 1, true, true)
	form.AddTextView("Location", fmt.Sprintf("%v %v %v", ui.country, ui.province, ui.city), 50, 1, true, true)
	form.AddTextView("ISP", ui.isp, 50, 1, true, true)
	form.AddTextView("Domain", ui.localDomain.Hostname(), 50, 1, true, true)
	form.AddTextView("DNS Lookup", ipText, 50, len(ip), true, true)
	//form.AddTextView("Nameserver", nsText, 50, len(ns), true, true)

	form.AddButton("Confirm", func() {
		//ui.appStop()

		// 下载 Caddy
		if !ui.supportProxyProtocol {
			ui.nextContent(ui.caddyBinaryDownload("xray", false))
		} else {
			ui.nextContent(ui.xrayInstall(false))
		}

	})

	form.AddButton("Back", func() {
		ui.nextContent(ui.tlsForm())
	})

	return ui.contentWithTitle("Confirm Information", form)
}

func (ui *UI) xrayVisionTLS() *UI {
	x := &XrayTLSConfig{}

	if err := json.Unmarshal([]byte(XrayTLSConfigExample), x); err != nil {
		panic(err)
	}

	sslCert, ok := ui.getCaddySSLCert()
	if !ok {
		panic("failed to get ssl cert")
	}

	sslKey, ok := ui.getCaddySSLKey()
	if !ok {
		panic("failed to get ssl key")
	}

	x.Inbounds[0].Port = ui.xrayPort
	x.Inbounds[0].Settings.Clients[0].Id = ui.xrayUUID
	x.Inbounds[0].Settings.Fallbacks[0].Dest = strconv.Itoa(ui.caddyProxyProtocolPort)
	x.Inbounds[0].StreamSettings.TlsSettings.Certificates[0].CertificateFile = sslCert
	x.Inbounds[0].StreamSettings.TlsSettings.Certificates[0].KeyFile = sslKey

	b, err := json.MarshalIndent(x, "", " ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("/usr/local/etc/xray/config.json", b, 644); err != nil {
		panic(err)
	}

	if _, err = command.XraySystemd("restart").CombinedOutput(); err != nil {
		panic(err)
	}

	ui.app.QueueUpdateDraw(func() {
		ui.updateFooter()
		ui.nextContent(ui.xrayQrcode())
	})

	return ui
}

const XrayTLSConfigExample = `
{
  "log": {
    "loglevel": "warning"
  },
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "domain": [
          "geosite:category-ads-all"
        ],
        "outboundTag": "block"
      },
      {
        "type": "field",
        "domain": [
          "geosite:google"
        ],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": [
          "geoip:cn"
        ],
        "outboundTag": "block"
      }
    ]
  },
  "inbounds": [
    {
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [
          {
            "id": "",
            "flow": "xtls-rprx-vision"
          }
        ],
        "decryption": "none",
        "fallbacks": [
          {
            "dest": "8080",
            "xver": 1
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "tls",
        "tlsSettings": {
          "rejectUnknownSni": true,
          "minVersion": "1.3",
          "certificates": [
            {
              "certificateFile": "",
              "keyFile": ""
            }
          ]
        }
      },
      "sniffing": {
        "enabled": true,
        "destOverride": [
          "http",
          "tls"
        ]
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "tag": "direct"
    },
    {
      "protocol": "blackhole",
      "tag": "block"
    }
  ],
  "policy": {
    "levels": {
      "0": {
        "handshake": 2,
        "connIdle": 120
      }
    }
  }
}`

type XrayTLSConfig struct {
	Log struct {
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			Domain      []string `json:"domain,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			Ip          []string `json:"ip,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	Inbounds []struct {
		Listen   string `json:"listen"`
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
		Settings struct {
			Clients []struct {
				Id   string `json:"id"`
				Flow string `json:"flow"`
			} `json:"clients"`
			Decryption string `json:"decryption"`
			Fallbacks  []struct {
				Dest string `json:"dest"`
				Xver int    `json:"xver"`
			} `json:"fallbacks"`
		} `json:"settings"`
		StreamSettings struct {
			Network     string `json:"network"`
			Security    string `json:"security"`
			TlsSettings struct {
				RejectUnknownSni bool   `json:"rejectUnknownSni"`
				MinVersion       string `json:"minVersion"`
				Certificates     []struct {
					CertificateFile string `json:"certificateFile"`
					KeyFile         string `json:"keyFile"`
				} `json:"certificates"`
			} `json:"tlsSettings"`
		} `json:"streamSettings"`
		Sniffing struct {
			Enabled      bool     `json:"enabled"`
			DestOverride []string `json:"destOverride"`
		} `json:"sniffing"`
	} `json:"inbounds"`
	Outbounds []struct {
		Protocol string `json:"protocol"`
		Tag      string `json:"tag"`
	} `json:"outbounds"`
	Policy struct {
		Levels struct {
			Field1 struct {
				Handshake int `json:"handshake"`
				ConnIdle  int `json:"connIdle"`
			} `json:"0"`
		} `json:"levels"`
	} `json:"policy"`
}
