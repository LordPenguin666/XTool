package main

import (
	"fmt"
	"go.uber.org/zap"
)

type Config struct {
	CaddyHTTPPort  int
	CaddyHTTPSPort int
	XrayXTLSPort   int

	Domain      string
	ProxyDomain string

	logger *zap.Logger

	Arch     string
	Platform string
	IP       string

	XrayVer string

	CaddyVer                  string
	CaddyProxyProtocolSupport bool
}

type Option func(c *Config)

func (opt Option) Apply(c *Config) {
	opt(c)
}

func DefaultOptions() *Config {
	c := &Config{
		CaddyHTTPPort:  8080,
		CaddyHTTPSPort: 8443,
		XrayXTLSPort:   443,
		logger:         logger(),
	}

	c = c.HostInfo().CaddyVersion().CaddyPlugins("caddy.listeners.proxy_protocol").
		XrayVersion().IPAddress()

	return c
}

func Ready(defaultOptions *Config, opts ...Option) *Config {
	c := defaultOptions

	for _, o := range opts {
		o.Apply(c)
	}

	c.DomainLookup()

	fmt.Println(c)

	return c
}

func WithXTLSPort(port int) Option {
	return func(c *Config) {
		c.XrayXTLSPort = port
	}
}

func WithCaddyHTTPPort(port int) Option {
	return func(c *Config) {
		c.CaddyHTTPPort = port
	}
}

func WithCaddyHTTPSPort(port int) Option {
	return func(c *Config) {
		c.CaddyHTTPSPort = port
	}
}

func WithDomain(domain string) Option {
	return func(c *Config) {
		c.Domain = domain
	}
}

func withProxyDomain(proxyDomain string) Option {
	return func(c *Config) {
		c.ProxyDomain = proxyDomain
	}
}
