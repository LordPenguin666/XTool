package main

import (
	"github.com/shirou/gopsutil/v3/host"
)

func (c *Config) HostInfo() *Config {
	h, err := host.Info()
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	c.Arch = h.KernelArch
	c.Platform = h.PlatformFamily

	return c
}
