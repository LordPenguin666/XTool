package main

import (
	"fmt"
	"os/exec"
)

func (c *Config) XrayInstallation() *Config {
	shell := exec.Command("bash", "-c",
		"$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)",
		"@ install -u caddy",
	)
	fmt.Println(shell.String())
	output, err := shell.CombinedOutput()
	if err != nil {
		c.logger.Error(err.Error())
	}
	fmt.Println(string(output))
	return c
}

func (c *Config) XrayUninstallation() *Config {
	shell := exec.Command("bash", "-c",
		"$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)",
		"@", "remove", "--purge",
	)
	output, err := shell.CombinedOutput()
	if err != nil {
		c.logger.Error(err.Error())
	}
	fmt.Println(string(output))
	return c
}
