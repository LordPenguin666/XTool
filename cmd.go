package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	opts := DefaultOptions()

	clear := exec.Command("clear")
	if err := clear.Run(); err != nil {
		opts.logger.Error(err.Error())
	}

	opts.menu()

	var input int

Loop:
	for {
		fmt.Printf("请输入数字 (回车确认): ")
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
	default:
		fmt.Printf("%v %v", red("[Warning]"), yellow("请输入一个正确的数字!\n\n"))
		goto Loop
	}
}
