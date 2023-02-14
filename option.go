package main

import (
	"go.uber.org/zap"
)

type Config struct {
	CaddyHTTPPort  int
	CaddyHTTPSPort int
	XrayXTLSPort   int

	Domain      string
	ProxyDomain string
	isHttps     bool

	logger *zap.Logger

	Arch     string
	Platform string
	IP       string

	XrayVer    string
	XrayLink   string
	XrayUUID   string
	XrayConfig *XrayConfig

	CaddyVer                  string
	CaddySSLCert              string
	CaddySSLKey               string
	CaddyProxyProtocolSupport bool

	PackageManagement string
}

//type Option func(c *Config)

//func (opt Option) Apply(c *Config) {
//	opt(c)
//}

func DefaultOptions() *Config {
	c := &Config{
		CaddyHTTPPort:  8080,
		CaddyHTTPSPort: 8443,
		XrayXTLSPort:   443,
		logger:         logger(),
	}

	c = c.HostInfo().SwitchPackage().CaddyVersion().
		CaddyPlugins("caddy.listeners.proxy_protocol").
		XrayVersion().IPAddress()

	return c
}

//func Ready(defaultOptions *Config, opts ...Option) *Config {
//	c := defaultOptions
//
//	for _, o := range opts {
//		o.Apply(c)
//	}
//
//	c.DomainLookup()
// //
//	return c
//}
//
//func WithXTLSPort(port int) Option {
//	return func(c *Config) {
//		c.XrayXTLSPort = port
//	}
//}
//
//func WithCaddyHTTPPort(port int) Option {
//	return func(c *Config) {
//		c.CaddyHTTPPort = port
//	}
//}
//
//func WithCaddyHTTPSPort(port int) Option {
//	return func(c *Config) {
//		c.CaddyHTTPSPort = port
//	}
//}
//
//func WithDomain(domain string) Option {
//	return func(c *Config) {
//		c.Domain = domain
//	}
//}
//
//func withProxyDomain(proxyDomain string) Option {
//	return func(c *Config) {
//		c.ProxyDomain = proxyDomain
//	}
//}
