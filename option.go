package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/host"
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

	CaddyVer                  string
	CaddyProxyProtocolSupport bool
}

type Option func(c *Config)

func (opt Option) Apply(c *Config) {
	opt(c)
}

func Ready(opts ...Option) *Config {
	c := &Config{
		CaddyHTTPPort:  8080,
		CaddyHTTPSPort: 8443,
		XrayXTLSPort:   443,
	}

	for _, o := range opts {
		o.Apply(c)
	}

	h, err := host.Info()
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	c.Arch = h.KernelArch
	c.Platform = h.Platform
	c.IP = c.IPAddress()
	ips := c.DomainLookup()

	var exist bool
	for _, ip := range ips {
		if ip.String() == c.IP {
			exist = true
			break
		}
	}

	if !exist {
		c.logger.Warn("please confirm your domain resolved is correct",
			data(
				"domain", c.Domain,
				"current_ip", c.IP,
				"lookup", ips)...)
	} else {
		color.Green(fmt.Sprintf("域名解析正确 %v -> %v", c.Domain, c.IP))
	}

	return c
}

func WithXTLSPort(port int) Option {
	return func(c *Config) {
		c.XrayXTLSPort = port
	}
}
