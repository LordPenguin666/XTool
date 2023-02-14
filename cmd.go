package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	opts := DefaultOptions()

	opts.Clear()
	opts.menu()

	var input int

Loop:
	for {
		fmt.Printf("%v%v: ", green("请输入数字"), blue(" (回车确认)"))
		if _, err := fmt.Scan(&input); err != nil {
			fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的数字!\n\n"))
			continue
		}
		break
	}

	switch input {
	case 0:
		os.Exit(0)
	case 1:
		opts = opts.ConfirmModify()
		opts.CaddyInstallation()
	case 2:

	case 11:
	case 12:
	case 13:
		opts.StartCaddy()
	case 14:
		opts.RestartCaddy()
	case 15:
		opts.StopCaddy()
	case 16:
		opts.EnableCaddy()
	case 17:
		opts.DisableCaddy()
	case 18:
		opts.UninstallCaddy()
	case 19:
		opts.CaddyInstallation()

	case 21:
	case 22:
	case 23:
	case 24:
	case 25:
	case 26:
	case 27:
	case 28:
	case 29:

	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的数字!\n\n"))
		goto Loop
	}
}

func (c *Config) Clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	if err := clear.Run(); err != nil {
		c.logger.Error(err.Error())
	}
}

func (c *Config) FileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}

	return true
}
