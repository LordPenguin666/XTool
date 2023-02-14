package main

import (
	"fmt"
	"os"
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

func (c *Config) CaddySSL() *Config {
	// debian 11
	var certExist, keyExist bool

	issuers := []string{
		"acme.zerossl.com-v2-dv90",
		"acme-v02.api.letsencrypt.org-directory",
	}

	for _, issuer := range issuers {
		cert := fmt.Sprintf("/var/lib/caddy/.local/share/caddy/certificates/%v/%v/%v.crt",
			issuer, c.Domain, c.Domain,
		)

		key := fmt.Sprintf("/var/lib/caddy/.local/share/caddy/certificates/%v/%v/%v.key",
			issuer, c.Domain, c.Domain,
		)

		if _, err := os.Stat(cert); err == nil {
			c.CaddySSLCert = cert
			certExist = true
		}

		if _, err := os.Stat(key); err == nil {
			c.CaddySSLCert = key
			keyExist = true
		}
	}

	if !certExist || !keyExist {
		fmt.Printf("%v %v", red("[Warning]"), yellow("没有找到 SSL 证书, 请确认域名是否正确解析, 对应的端口是否开放, 在确认无误后请提交 issue!\n\n"))
		os.Exit(0)
	}

	return c
}
