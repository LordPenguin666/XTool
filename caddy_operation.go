package main

import "os/exec"

func (c *Config) StartCaddy() *Config {
	bash := exec.Command("systemctl", "start", "caddy")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) RestartCaddy() *Config {
	bash := exec.Command("systemctl", "restart", "caddy")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) StopCaddy() *Config {
	bash := exec.Command("systemctl", "stop", "caddy")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) EnableCaddy() *Config {
	bash := exec.Command("systemctl", "enable", "caddy")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) DisableCaddy() *Config {
	bash := exec.Command("systemctl", "disable", "caddy")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) UninstallCaddy() *Config {
	return c
}
