package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"os"
	"os/exec"
	"strings"
)

func (c *Config) HostInfo() *Config {
	h, err := host.Info()
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	//fmt.Println(h.KernelArch, h.Platform, h.PlatformFamily, h.PlatformVersion)

	c.Arch = h.KernelArch
	c.Platform = h.PlatformFamily

	return c
}

func (c *Config) SwitchPackage() *Config {
	switch c.Platform {
	case "debian":
		c.PackageManagement = "apt"
	case "rhel":
		whereIsYum := exec.Command("whereis", "yum")
		whereIsDnf := exec.Command("whereis", "dnf")

		output, err := whereIsYum.CombinedOutput()
		if err != nil {
			c.logger.Error(err.Error())
		}

		words := strings.Split(string(output), " ")
		if len(words) <= 1 {
			output, err = whereIsDnf.CombinedOutput()
			if err != nil {
				c.logger.Error(err.Error())
			}

			words = strings.Split(string(output), " ")

			if len(words) <= 1 {
				fmt.Printf("%v", red("?????????????????????????"))
				os.Exit(0)
			}

			c.PackageManagement = "dnf"

		} else {
			c.PackageManagement = "yum"
		}
	}
	return c
}
