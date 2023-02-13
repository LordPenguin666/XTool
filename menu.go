package main

import "github.com/fatih/color"

func (c *Config) menu() {
	color.Blue("*********************************")
	color.Blue("*                               *")

	if _, err := color.New(color.FgBlue).Print("*"); err != nil {
		c.logger.Error(err.Error())
	}

	if _, err := color.New(color.FgRed).Print("        Xray tool v1.0         "); err != nil {
		c.logger.Error(err.Error())

	}

	if _, err := color.New(color.FgBlue).Println("*"); err != nil {
		c.logger.Error(err.Error())
	}

	color.Blue("*         (github)              *")
}
