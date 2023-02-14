package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"os/exec"
)

const XrayScript = "install-release.sh"

func (c *Config) XrayInstallation() *Config {
	// exec好像有点问题, 只能分步了
	fmt.Printf("%v\n", green("下载 Xray 脚本中..."))
	client := resty.New()
	client.SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	if _, err := client.R().SetOutput(XrayScript).
		Get("https://github.com/XTLS/Xray-install/raw/main/install-release.sh"); err != nil {
		c.logger.Error(err.Error())
	}

	if err := os.Chmod(XrayScript, 0755); err != nil {
		c.logger.Error(err.Error())
	}

	fmt.Printf("%v\n", green("执行 Xray 安装中..."))
	install := exec.Command("bash", "-c", "./install-release.sh install -u caddy")
	output, err := install.CombinedOutput()
	if err != nil {
		c.logger.Error(err.Error())
	}

	fmt.Println(string(output))

	if err = os.Remove("install-release.sh"); err != nil {
		c.logger.Error(err.Error())
	}

	// 安装 Xray
	c.LoadXrayConfig().RestartXray().EnableXray()

	return c
}

func (c *Config) XrayUninstallation() *Config {
	fmt.Printf("%v\n", green("下载 Xray 脚本中..."))
	client := resty.New()
	client.SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	if _, err := client.R().SetOutput("install-release.sh").
		Get("https://github.com/XTLS/Xray-install/raw/main/install-release.sh"); err != nil {
		c.logger.Error(err.Error())
	}

	if err := os.Chmod(XrayScript, 0755); err != nil {
		c.logger.Error(err.Error())
	}

	fmt.Printf("%v\n", green("卸载 Xray 中..."))
	install := exec.Command("bash", "-c", "./install-release.sh remove --purge")
	output, err := install.CombinedOutput()
	if err != nil {
		c.logger.Error(err.Error())
	}

	fmt.Println(string(output))

	if err = os.Remove("install-release.sh"); err != nil {
		c.logger.Error(err.Error())
	}

	return c
}
