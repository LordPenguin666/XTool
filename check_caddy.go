package main

import (
	"os/exec"
	"strings"
)

func (c *Config) CaddyVersion() *Config {
	cmd := exec.Command("caddy", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c
	}

	//c.logger.Debug("", data("out", string(out))...)

	words := strings.Split(string(out), " ")

	if len(words) == 0 {
		return c
	}

	c.CaddyVer = words[0]

	return c
}

func (c *Config) CaddyPlugins(name string) *Config {
	cmd := exec.Command("caddy", "list-modules")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c
	}

	if !strings.Contains(string(out), name) {
		return c
	}

	c.CaddyProxyProtocolSupport = true
	return c
}
