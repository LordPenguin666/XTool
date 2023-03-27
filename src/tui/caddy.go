package tui

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/rivo/tview"
	"net/url"
	"os"
	"os/exec"
	"time"
	"xray-tool/src/command"
)

func (ui *UI) caddyMenu() *tview.Grid {
	list := tview.NewList()

	list.AddItem("Start caddy", "Caddy 启动", '0', func() {
		_, err := command.CaddySystemd("start").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Restart caddy", "Caddy 重启", '1', func() {
		_, err := command.CaddySystemd("restart").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Stop caddy", "Caddy 关闭", '2', func() {
		_, err := command.CaddySystemd("stop").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Enable caddy", "Caddy 开机启动", '3', func() {
		_, err := command.CaddySystemd("enable").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Disable caddy", "Caddy 关闭开机启动", '4', func() {
		_, err := command.CaddySystemd("disable").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Install caddy", "Caddy 安装 / 更新", '5', func() {
		ui.nextContent(ui.caddyBinaryDownload("menu", true)).updateFooter()
	})

	list.AddItem("Uninstall caddy", "Caddy 卸载", '6', func() {
		switch ui.platform {
		case "rhel":
			if output, err := command.CaddyRedHatUninstall().CombinedOutput(); err != nil {
				panic(string(output))
			}
		case "debian":
			if output, err := command.CaddyDebianUninstall().CombinedOutput(); err != nil {
				panic(string(output))
			}
		default:
			panic(fmt.Sprintf("not a support platform system: %v", ui.platform))
		}

		if localDomain, ok := ui.getLocalDomain(); ok {
			ui.localDomain = localDomain
		} else {
			ui.localDomain, _ = url.Parse("")
		}

		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Back", "返回", 'c', func() {
		ui.nextContent(ui.menu())
	})

	return ui.contentWithTitle("Caddy Menu", list)
}

func (ui *UI) caddyBinaryDownload(next string, update bool) *tview.Grid {
	text := "Start to Download Caddy...\n"
	textView := tview.NewTextView().SetTextAlign(tview.AlignLeft)

	go func() {
		//time.Sleep(time.Second)

		ifFileExist := func() {
			// 如果 caddy 存在, 则删除
			if ui.fileExist("/usr/bin/caddy") {
				if _, err := command.CaddySystemd("stop").CombinedOutput(); err != nil {
					panic(err)
				}

				if err := os.Remove("/usr/bin/caddy"); err != nil {
					panic(err)
				}

				ui.app.QueueUpdateDraw(func() {
					ui.updateFooter()
				})
			}
		}
		u := url.Values{}

		u.Add("os", "linux")
		//u.Add("p", "github.com/mholt/caddy-l4")
		u.Add("p", "github.com/mastercactapus/caddy2-proxyprotocol")
		//u.Add("p", "github.com/caddyserver/jsonc-adapter")

		switch ui.arch {
		case "x86_64":
			u.Add("arch", "amd64")
		case "aarch64":
			u.Add("arch", "arm64")
		default:
			u.Add("arch", ui.arch)
		}

		ifFileExist()

		// 从官方下载 Caddy 配置
		var cmds []*exec.Cmd

		switch ui.platform {
		case "rhel":
			cmds = command.CaddyRedHatInstall()
		case "debian":
			cmds = command.CaddyDebianInstall()
		//case "arch":
		//	cmds = command.CaddyArchInstall()
		default:
			panic("not a supported system: " + ui.platform)
		}

		for _, cmd := range cmds {
			output, _ := cmd.CombinedOutput()
			ui.app.QueueUpdateDraw(func() {
				text += "\n" + string(output)
			})
		}

		ui.updateFooter()
		ifFileExist()

		// 替换为有插件的 Caddy
		now := time.Now()
		callback := func(info req.DownloadInfo) {
			if info.Response.Response != nil {
				text += fmt.Sprintf("\nDownloading: %v", time.Since(now))
				ui.app.QueueUpdateDraw(func() {
					textView.SetText(text)
				})
			}
		}

		_, err := ui.client.R().
			SetQueryString(u.Encode()).
			SetOutputFile("/usr/bin/caddy").
			SetDownloadCallback(callback).
			Get("https://caddyserver.com/api/download")

		if err != nil {
			panic(err)
		}

		if err = os.Chmod("/usr/bin/caddy", 755); err != nil {
			panic(err)
		}

		if update {
			if output, err := command.CaddySystemd("restart").CombinedOutput(); err != nil {
				panic(string(output))
			}
			return
		}

		// 配置 /etc/caddy/conf.d
		if _, err := command.CaddyConfigMkdir().CombinedOutput(); err != nil {
			panic(err)
		}

		// 配置 Caddyfile
		if err = os.WriteFile("/etc/caddy/Caddyfile", []byte(ui.getMainCaddyFile()), 644); err != nil {
			panic(err)
		}

		// 配置 站点
		switch ui.localDomain.Scheme {
		case "http":
			if err = os.WriteFile(fmt.Sprintf("/etc/caddy/conf.d/%v", ui.localDomain.Hostname()), []byte(ui.getCaddyfileProxyHTTP()), 644); err != nil {
				panic(err)
			}
		case "https":
			if err = os.WriteFile(fmt.Sprintf("/etc/caddy/conf.d/%v", ui.localDomain.Hostname()), []byte(ui.getCaddyfileProxyHTTPS()), 644); err != nil {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("not a support URI: %v", ui.localDomain))
		}

		// 重载 Systemd
		//if output, _ := command.SystemdReload().CombinedOutput(); err != nil {
		//	panic(string(output))
		//}

		// 重启 Caddy
		if output, err := command.CaddySystemd("restart").CombinedOutput(); err != nil {
			panic(string(output))
		}

		if output, err := command.CaddySystemd("enable").CombinedOutput(); err != nil {
			panic(string(output))
		}

		_, certExist := ui.getCaddySSLCert()
		_, keyExist := ui.getCaddySSLKey()

		if !certExist || !keyExist {
			text += "\n\nWait 30 seconds for applying the certification"

			ui.app.QueueUpdateDraw(func() {
				textView.SetText(text)
			})

			time.Sleep(30 * time.Second)
		}

		switch next {
		case "menu":
			ui.nextContent(ui.menu())
		case "xray":
			ui.nextContent(ui.xrayInstall(false))
		}
	}()

	return ui.contentWithTitle("Caddy Installing", textView)
}

func (ui *UI) getPort(u *url.URL) string {
	if u.Port() == "" {
		if u.Scheme == "http" {
			return "80"
		} else if u.Scheme == "https" {
			return "443"
		} else {
			panic("if you are using h2c, please add port to your uri.")
		}
	}

	return u.Port()
}
