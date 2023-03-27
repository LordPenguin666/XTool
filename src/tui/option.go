package tui

import (
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/host"
	"net/url"
	"strings"
)

type UI struct {
	app     *tview.Application
	layout  *tview.Grid
	header  *tview.Grid
	content *tview.Grid
	footer  *tview.Grid

	// system
	arch                 string
	platform             string
	xrayInstalled        bool
	supportProxyProtocol bool
	//supportLayer4 bool

	// network
	isp      string
	addr     string
	city     string
	country  string
	province string

	// config
	xrayUUID               string
	xrayPort               int
	caddyHTTPSPort         int
	caddyProxyProtocolPort int

	client      *req.Client
	xrayFunc    func() *UI
	proxyURI    *url.URL
	localDomain *url.URL
}

func New() *UI {
	ui := &UI{
		app: tview.NewApplication(),

		arch:     getArch(),
		platform: getPlatform(),

		xrayUUID:               uuid.NewString(),
		xrayPort:               443,
		caddyProxyProtocolPort: 8080,
		caddyHTTPSPort:         8443,

		client:      req.NewClient().SetUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
		localDomain: getDefaultURI(),
	}

	if localDomain, ok := ui.getLocalDomain(); ok {
		ui.localDomain = localDomain
	}

	ui.newLayout().newHeader().newContent().newFooter()
	ui.updateContent().updateLayout().getIP().appStart()

	return ui
}

func getArch() string {
	h, err := host.Info()
	if err != nil {
		panic(err)
	}

	return h.KernelArch
}

func getPlatform() string {
	h, err := host.Info()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(h.PlatformFamily)
}

func getDefaultURI() *url.URL {
	u, err := url.Parse("https://")
	if err != nil {
		panic(err)
	}

	return u
}

func (ui *UI) getIP() *UI {

	zoneResp := &ZoneResp{}

	if _, err := ui.client.R().SetSuccessResult(zoneResp).Get("https://api.bilibili.com/x/web-interface/zone"); err != nil {
		panic(err)
	}

	ui.isp = zoneResp.Data.Isp
	ui.addr = zoneResp.Data.Addr
	ui.city = zoneResp.Data.City
	ui.country = zoneResp.Data.Country
	ui.province = zoneResp.Data.Province

	return ui
}

func (ui *UI) appStart() *UI {
	if err := ui.app.SetRoot(ui.layout, true).SetFocus(ui.content).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return ui
}

func (ui *UI) appStop() *UI {
	ui.app.Stop()

	return ui
}

type ZoneResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Addr        string  `json:"addr"`
		Country     string  `json:"country"`
		Province    string  `json:"province"`
		City        string  `json:"city"`
		Isp         string  `json:"isp"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		ZoneId      int     `json:"zone_id"`
		CountryCode int     `json:"country_code"`
	} `json:"data"`
}
