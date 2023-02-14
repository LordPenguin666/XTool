package main

import "os/exec"

func (c *Config) StartXray() *Config {
	bash := exec.Command("systemctl", "start", "xray")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) RestartXray() *Config {
	bash := exec.Command("systemctl", "restart", "xray")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) StopXray() *Config {
	bash := exec.Command("systemctl", "stop", "xray")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) EnableXray() *Config {
	bash := exec.Command("systemctl", "enable", "xray")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}

func (c *Config) DisableXray() *Config {
	bash := exec.Command("systemctl", "disable", "xray")
	if err := bash.Run(); err != nil {
		logger().Error(err.Error())
	}
	return c
}
