package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
)

func (c *Config) HostInfo() *Config {
	h, err := host.Info()
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	fmt.Println(h.Platform, h.PlatformFamily, h.PlatformVersion)

	c.Arch = h.KernelArch
	c.Platform = h.Platform

	return c
}
