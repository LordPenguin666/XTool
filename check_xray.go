package main

import (
	"os/exec"
	"strings"
)

func (c *Config) XrayVersion() *Config {
	cmd := exec.Command("xray", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c
	}

	//c.logger.Debug("", data("out", string(out))...)
	words := strings.Split(string(out), " ")

	if len(words) == 0 {
		return c
	}

	c.XrayVer = words[1]

	return c
}
