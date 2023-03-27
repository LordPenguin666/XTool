package tui

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/headers"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	caddy2proxyprotocol "github.com/mastercactapus/caddy2-proxyprotocol"
	"github.com/mholt/caddy-l4/layer4"
	"github.com/mholt/caddy-l4/modules/l4proxy"
	"net/http"
	"net/url"
	"os"
)

func (ui *UI) getLocalDomain() (*url.URL, bool) {
	b, err := os.ReadFile("/etc/caddy/Caddyfile")
	if err != nil {
		return nil, false
	}

	outputs, err := caddyfile.Parse("/etc/caddy/Caddyfile", b)
	if err != nil {
		return nil, false
	}

	for _, output := range outputs {
		for _, key := range output.Keys {
			u, err := url.Parse(key)
			if err != nil {
				continue
			}
			if u.Hostname() != "" {
				return u, true
			}
		}
	}

	return nil, false
}

func (ui *UI) getMainCaddyFile() string {
	return "{\n" +
		fmt.Sprintf(
			"        servers :%v {\n", ui.caddyProxyProtocolPort) +
		"                listener_wrappers {\n" +
		"                        proxy_protocol {\n" +
		"                                timeout 2s\n" +
		"                                allow 0.0.0.0/0\n" +
		"                        }\n" +
		"                        tls\n" +
		"                }\n" +
		"                protocols h1 h2 h2c h3\n" +
		"        }\n" +
		"}\n" +
		"\n" +
		":80 {\n" +
		"    redir https://{host}{url}\n" +
		"}\n" +
		"\n" +
		"import /etc/caddy/conf.d/*"
}

func (ui *UI) getCaddyfileProxyHTTP() string {
	return fmt.Sprintf("http://%v:%v {\n", ui.localDomain.Hostname(), ui.caddyProxyProtocolPort) +
		fmt.Sprintf("    reverse_proxy %v\n", ui.localDomain) +
		"}\n" +
		"\n" +
		fmt.Sprintf("https://%v:%v {\n", ui.localDomain.Hostname(), ui.caddyHTTPSPort) +
		fmt.Sprintf("    reverse_proxy %v\n", ui.localDomain) +
		"}"
}

func (ui *UI) getCaddyfileProxyHTTPS() string {
	return fmt.Sprintf("http://%v:%v {\n", ui.localDomain.Hostname(), ui.caddyProxyProtocolPort) +
		fmt.Sprintf("    reverse_proxy %v {\n", ui.proxyURI) +
		"        header_up Host {upstream_hostport}\n" +
		"        transport http {\n" +
		"            tls\n" +
		"        }\n" +
		"    }\n" +
		"}\n" +
		"\n" +
		fmt.Sprintf("%v:%v {\n", ui.localDomain.Hostname(), ui.caddyHTTPSPort) +
		fmt.Sprintf("    reverse_proxy %v {\n", ui.proxyURI) +
		"        header_up Host {upstream_hostport}\n" +
		"        transport http {\n" +
		"            tls\n" +
		"        }\n" +
		"    }\n" +
		"}"
}

func (ui *UI) CaddyConfigWithLayer4() *UI {
	config := &caddy.Config{
		Admin: &caddy.AdminConfig{Disabled: true},
		AppsRaw: map[string]json.RawMessage{
			"http": caddyconfig.JSON(&caddyhttp.App{
				Servers: map[string]*caddyhttp.Server{
					ui.localDomain.Hostname(): {
						Protocols: []string{"h1", "h2", "h2c", "h3"},
						Listen:    []string{fmt.Sprintf(":%v", ui.caddyProxyProtocolPort)},
						AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
							Disabled: true,
						},
						ListenerWrappersRaw: []json.RawMessage{
							caddyconfig.JSONModuleObject(&caddy2proxyprotocol.Wrapper{
								Allow: []string{"0.0.0.0/0"},
							}, "wrapper", "proxy_protocol", nil),
							caddyconfig.JSONModuleObject(&caddy2proxyprotocol.Wrapper{}, "wrapper", "tls", nil),
						},
						Routes: []caddyhttp.Route{{
							MatcherSetsRaw: []caddy.ModuleMap{{
								"host": caddyconfig.JSON(caddyhttp.MatchHost{ui.localDomain.Hostname()}, nil),
							}},
							HandlersRaw: []json.RawMessage{
								caddyconfig.JSONModuleObject(&reverseproxy.Handler{
									Upstreams: []*reverseproxy.Upstream{{
										Dial: fmt.Sprintf("%v:%v", ui.proxyURI.Hostname(), ui.getPort(ui.proxyURI)),
									}},
									Headers: &headers.Handler{
										Request: &headers.HeaderOps{
											Set: http.Header{
												"Host": []string{ui.proxyURI.Hostname()},
											},
										},
									},
									TransportRaw: caddyconfig.JSONModuleObject(&reverseproxy.HTTPTransport{
										TLS: func() *reverseproxy.TLSConfig {
											switch ui.proxyURI.Scheme {
											case "https":
												return &reverseproxy.TLSConfig{}
											default:
												return nil
											}
										}(),
									}, "protocol", "http", nil),
								}, "handler", "reverse_proxy", nil),
							},
						}},
					},
				},
			}, nil),
			"layer4": caddyconfig.JSON(
				&layer4.App{Servers: map[string]*layer4.Server{
					ui.localDomain.Hostname(): {
						Listen: []string{":443"},
						Routes: []*layer4.Route{{
							MatcherSetsRaw: []caddy.ModuleMap{{
								"tls": caddyconfig.JSON(caddy.ModuleMap{
									"sni": caddyconfig.JSON([]string{ui.localDomain.Hostname()}, nil),
								}, nil),
							}},
							HandlersRaw: []json.RawMessage{
								caddyconfig.JSONModuleObject(l4proxy.Handler{
									Upstreams: []*l4proxy.Upstream{{
										Dial: []string{
											fmt.Sprintf("%v:%v", ui.localDomain.Hostname(), ui.xrayPort),
										},
									}},
								}, "handler", "proxy", nil),
							},
						}},
					},
				}}, nil),
			"tls": caddyconfig.JSON(
				&caddytls.TLS{
					Automation: &caddytls.AutomationConfig{
						Policies: []*caddytls.AutomationPolicy{{
							Subjects: []string{ui.localDomain.Hostname()},
							//Issuers: []certmagic.Issuer{},
						}},
					},
				}, nil,
			),
		},
	}

	b, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("/etc/caddy/config.json", b, 644); err != nil {
		panic(err)
	}

	return ui
}
