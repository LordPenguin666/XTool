package main

import "os/exec"

func (c *Config) CaddyVersion() string {
	cmd := exec.Command("caddy", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	out
}
