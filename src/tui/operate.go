package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
	"xray-tool/src/command"
)

func (ui *UI) getCaddyVersion() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Caddy: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)

	output, err := command.CaddyVersion().CombinedOutput()
	if err != nil {
		grid.AddItem(tview.NewTextView().SetText("未安装").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
		return grid
	}

	version := strings.Split(string(output), " ")[0]
	grid.AddItem(tview.NewTextView().SetText(fmt.Sprintf("%v", version)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)

	return grid
}

func (ui *UI) getProxyProtocolSupport() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Proxy Protocol: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)

	output, err := command.CaddyListModules().CombinedOutput()
	if err != nil {
		grid.AddItem(tview.NewTextView().SetText("unsupported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
		return grid
	}

	if !strings.Contains(string(output), "caddy.listeners.proxy_protocol") {
		grid.AddItem(tview.NewTextView().SetText("unsupported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
		return grid
	}

	ui.supportProxyProtocol = true
	grid.AddItem(tview.NewTextView().SetText("supported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)

	return grid
}

//func (ui *UI) getLayer4Support() *tview.Grid {
//	grid := tview.NewGrid().SetColumns(0, 0)
//	grid.AddItem(tview.NewTextView().SetText("Layer4: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)
//
//	output, err := command.CaddyListModules().CombinedOutput()
//	if err != nil {
//		grid.AddItem(tview.NewTextView().SetText("unsupported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
//		return grid
//	}
//
//	if !strings.Contains(string(output), "caddy.listeners.layer4") {
//		grid.AddItem(tview.NewTextView().SetText("unsupported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
//		return grid
//	}
//
//	ui.supportLayer4 = true
//	grid.AddItem(tview.NewTextView().SetText("supported").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)
//
//	return grid
//}

func (ui *UI) getXrayVersion() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Xray: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)

	output, err := command.XrayVersion().CombinedOutput()
	if err != nil {
		grid.AddItem(tview.NewTextView().SetText("未安装").SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
		return grid
	}

	version := strings.Split(string(output), " ")[1]
	grid.AddItem(tview.NewTextView().SetText(fmt.Sprintf("v%v", version)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)

	ui.xrayInstalled = true

	return grid
}

func (ui *UI) getPlatform() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Platform Family: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(ui.platform).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)

	return grid
}

//func (ui *UI) getArch() *tview.Grid {
//	grid := tview.NewGrid().SetColumns(0, 0)
//	grid.AddItem(tview.NewTextView().SetText("Arch: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)
//	grid.AddItem(tview.NewTextView().SetText(ui.arch).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)
//
//	return grid
//}

func (ui *UI) getCaddyState() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Caddy State: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)

	output, _ := command.CaddySystemd("is-active").CombinedOutput()
	//if err != nil {
	//	panic(err)
	//}

	state := strings.TrimSpace(string(output))

	switch state {
	case "active":
		grid.AddItem(tview.NewTextView().SetText(string(output)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)
	default:
		grid.AddItem(tview.NewTextView().SetText(string(output)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
	}

	return grid
}

func (ui *UI) getXrayState() *tview.Grid {
	grid := tview.NewGrid().SetColumns(0, 0)
	grid.AddItem(tview.NewTextView().SetText("Xray State: ").SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorLightBlue), 0, 0, 1, 1, 0, 0, false)

	output, _ := command.XraySystemd("is-active").CombinedOutput()
	//if err != nil {
	//	panic(err)
	//}

	state := strings.TrimSpace(string(output))

	switch state {
	case "active":
		grid.AddItem(tview.NewTextView().SetText(string(output)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorGreen), 0, 1, 1, 1, 0, 0, false)
	default:
		grid.AddItem(tview.NewTextView().SetText(string(output)).SetTextAlign(tview.AlignLeft).SetTextColor(tcell.ColorRed), 0, 1, 1, 1, 0, 0, false)
	}

	return grid
}

func (ui *UI) fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (ui *UI) getCaddySSLKey() (string, bool) {
	output, _ := command.CaddyKeyPath(command.CaddyCertPath(), ui.localDomain.Hostname(), "key").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	output, _ = command.CaddyKeyPath(command.CaddyRootCertPath(), ui.localDomain.Hostname(), "key").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	output, _ = command.CaddyKeyPath(command.CaddyShareCertPath(), ui.localDomain.Hostname(), "key").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	return "", false
}

func (ui *UI) getCaddySSLCert() (string, bool) {
	output, _ := command.CaddyKeyPath(command.CaddyCertPath(), ui.localDomain.Hostname(), "crt").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	output, _ = command.CaddyKeyPath(command.CaddyShareCertPath(), ui.localDomain.Hostname(), "crt").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	output, _ = command.CaddyKeyPath(command.CaddyRootCertPath(), ui.localDomain.Hostname(), "crt").CombinedOutput()

	if strings.Contains(string(output), ui.localDomain.Hostname()) {
		return strings.TrimSpace(string(output)), true
	}

	return "", false
}
