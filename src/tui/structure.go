package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) newLayout() *UI {
	ui.layout = tview.NewGrid().SetRows(5, 0, 2).SetBorders(true)

	return ui
}

func (ui *UI) newHeader() *UI {
	ui.header = tview.NewGrid().SetRows(0, 0, 0, 0, 0).
		AddItem(tview.NewTextView(), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Xray Tool v1.1").SetTextColor(tcell.ColorLightGreen), 1, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("https://github.com/LordPenguin666/XTool").SetTextColor(tcell.ColorLightBlue), 2, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("给面包狗点个 Star 喵, 给面包狗点个 Star 谢谢喵").SetTextColor(tcell.ColorOrangeRed), 3, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView(), 4, 0, 1, 1, 0, 0, false)

	return ui
}

func (ui *UI) newContent() *UI {
	ui.content = tview.NewGrid()

	return ui
}

func (ui *UI) newFooter() *UI {
	caddyInfo := tview.NewGrid().SetColumns(0, 0, 0)
	caddyInfo.AddItem(ui.getCaddyVersion(), 0, 0, 1, 1, 0, 0, false)
	caddyInfo.AddItem(ui.getCaddyState(), 0, 1, 1, 1, 0, 0, false)
	caddyInfo.AddItem(ui.getProxyProtocolSupport(), 0, 2, 1, 1, 0, 0, false)
	//caddyInfo.AddItem(ui.getLayer4Support(), 0, 3, 1, 1, 0, 0, false)

	xrayInfo := tview.NewGrid().SetColumns(0, 0, 0)
	xrayInfo.AddItem(ui.getXrayVersion(), 0, 0, 1, 1, 0, 0, false)
	xrayInfo.AddItem(ui.getXrayState(), 0, 1, 1, 1, 0, 0, false)
	xrayInfo.AddItem(ui.getPlatform(), 0, 2, 1, 1, 0, 0, false)
	//xrayInfo.AddItem(ui.getArch(), 0, 3, 1, 1, 0, 0, false)

	ui.footer = tview.NewGrid().SetRows(0, 0).
		AddItem(caddyInfo, 0, 0, 1, 1, 0, 0, false).
		AddItem(xrayInfo, 1, 0, 1, 1, 0, 0, false)

	return ui
}

func (ui *UI) updateLayout() *UI {
	ui.layout.AddItem(ui.header, 0, 0, 1, 1, 0, 0, false)
	ui.layout.AddItem(ui.content, 1, 0, 1, 1, 0, 0, true)
	ui.layout.AddItem(ui.footer, 2, 0, 1, 1, 0, 0, false)

	return ui
}

//func (ui *UI) updateHeader() *UI {
//
//	return ui
//}

func (ui *UI) updateContent() *UI {
	ui.content = ui.contentLayout()
	ui.content.AddItem(ui.menu(), 1, 1, 1, 1, 0, 0, true)

	return ui
}

func (ui *UI) updateFooter() *UI {
	ui.layout.RemoveItem(ui.footer)
	ui.newFooter()
	ui.layout.AddItem(ui.footer, 2, 0, 1, 1, 0, 0, false)

	return ui
}

func (ui *UI) contentWithTitle(title string, content tview.Primitive) *tview.Grid {
	return tview.NewGrid().SetBorders(true).SetRows(1, 0).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(title), 0, 0, 1, 1, 0, 0, false).
		AddItem(content, 1, 0, 1, 1, 0, 0, true)
}

func (ui *UI) nextContent(content *tview.Grid) *UI {
	ui.layout.RemoveItem(ui.content)

	ui.content = ui.contentLayout()
	ui.content.AddItem(content, 1, 1, 1, 1, 0, 0, true)

	ui.layout.AddItem(ui.content, 1, 0, 1, 1, 0, 0, true)

	// 重新定位输入框
	ui.app.SetFocus(ui.content)

	return ui
}

func (ui *UI) contentLayout() *tview.Grid {
	grid := tview.NewGrid().SetRows(1, 0).SetColumns(20, 0, 20)

	grid.AddItem(tview.NewTextView(), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView(), 1, 2, 1, 1, 0, 0, false)

	return grid
}
