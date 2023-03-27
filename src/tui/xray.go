package tui

import (
	"github.com/rivo/tview"
	"os"
	"xray-tool/src/command"
)

func (ui *UI) xrayMenu() *tview.Grid {
	list := tview.NewList()

	list.AddItem("Start xray", "Xray 启动", '0', func() {
		_, err := command.XraySystemd("start").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Restart xray", "Xray 重启", '1', func() {
		_, err := command.XraySystemd("restart").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Stop xray", "Xray 关闭", '2', func() {
		_, err := command.XraySystemd("stop").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Enable xray", "Xray 开机启动", '3', func() {
		_, err := command.XraySystemd("enable").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Disable xray", "Xray 关闭开机启动", '4', func() {
		_, err := command.XraySystemd("disable").CombinedOutput()
		if err != nil {
			panic(err)
		}
		ui.nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Install xray", "Xray 安装 / 更新", '5', func() {
		ui.nextContent(ui.xrayInstall(true))
	})

	list.AddItem("Uninstall xray", "Xray 卸载", '6', func() {
		ui.xrayScriptDownload()

		if _, err := command.XrayUninstall().CombinedOutput(); err != nil {
			panic(err)
		}

		ui.xrayScriptRemove().nextContent(ui.menu()).updateFooter()
	})

	list.AddItem("Back", "返回", 'c', func() {
		ui.nextContent(ui.menu())
	})

	return ui.contentWithTitle("Xray Menu", list)
}

func (ui *UI) xrayScriptDownload() *UI {
	_, err := ui.client.R().
		SetOutputFile("install-release.sh").
		Get("https://github.com/XTLS/Xray-install/raw/main/install-release.sh")

	if err != nil {
		panic(err)
	}

	if err = os.Chmod("install-release.sh", 0755); err != nil {
		panic(err)
	}

	return ui
}

func (ui *UI) xrayScriptRemove() *UI {
	if err := os.Remove("./install-release.sh"); err != nil {
		panic(err)
	}

	return ui
}

func (ui *UI) xrayInstall(update bool) *tview.Grid {
	text := "Start to install Xray...\n"
	textView := tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText(text)

	go func() {
		ui.xrayScriptDownload()

		text += "Xray install script download success\n\n"
		ui.app.QueueUpdateDraw(func() {
			textView.SetText(text)
		})

		output, err := command.XrayInstall().CombinedOutput()
		if err != nil {
			panic(err)
		}

		text += string(output)
		ui.app.QueueUpdateDraw(func() {
			textView.SetText(text)
		})

		ui.xrayScriptRemove()

		ui.app.QueueUpdateDraw(func() {
			ui.updateFooter()
		})

		if update {
			ui.nextContent(ui.menu())
			return
		}

		ui.xrayFunc()
	}()

	return ui.contentWithTitle("Xray Installation", textView)
}
