package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"os"
)

func (c *Config) LoadXrayConfig() *Config {
	if err := json.Unmarshal([]byte(XrayExampleConfig), c.XrayConfig); err != nil {
		c.logger.Error(err.Error())
	}

	c.XrayUUID = uuid.NewString()
	c.XrayConfig.Inbounds[0].Settings.Clients[0].Id = c.XrayUUID
	c.XrayConfig.Inbounds[0].Settings.Fallbacks[0].Dest = fmt.Sprintf("%v", c.CaddyHTTPPort)
	c.XrayConfig.Inbounds[0].StreamSettings.TlsSettings.Certificates[0].CertificateFile = c.CaddySSLCert
	c.XrayConfig.Inbounds[0].StreamSettings.TlsSettings.Certificates[0].KeyFile = c.CaddySSLKey

	b, err := json.Marshal(c.XrayConfig)
	if err != nil {
		c.logger.Error(err.Error())
	}

	if err = os.WriteFile("/usr/local/etc/xray/config.json", b, 644); err != nil {
		c.logger.Error(err.Error())
	}

	return c
}

func (c *Config) XrayLinkGeneration() *Config {
	u := url.Values{}
	u.Add("security", "tls")
	u.Add("encryption", "none")
	u.Add("alpn", "h2,http/1.1")
	u.Add("headerType", "none")
	u.Add("fp", "randomized")
	u.Add("type", "tcp")
	u.Add("flow", "xtls-rprx-vision")
	u.Add("sni", fmt.Sprintf("%v#[%v] (%v) xtls-rprx-vision", c.Domain, c.Domain, c.IP))

	c.XrayLink = fmt.Sprintf("vless://%v@%v:%v?%v",
		c.XrayUUID, c.Domain, c.XrayXTLSPort, u.Encode())

	return c
}

func (c *Config) XrayLinkPrint() *Config {
	fmt.Printf("%v %v",
		green("vless 链接:"),
		blue(c.XrayLink),
	)
	return c
}

const XrayExampleConfig = `
{
  "log": {
    "loglevel": "warning"
  },
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "domain": [
          "geosite:category-ads-all"
        ],
        "outboundTag": "block"
      },
      {
        "type": "field",
        "domain": [
          "geosite:google"
        ],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": [
          "geoip:cn"
        ],
        "outboundTag": "block"
      }
    ]
  },
  "inbounds": [
    {
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [
          {
            "id": "",
            "flow": "xtls-rprx-vision"
          }
        ],
        "decryption": "none",
        "fallbacks": [
          {
            "dest": "8080",
            "xver": 1
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "tls",
        "tlsSettings": {
          "rejectUnknownSni": true,
          "minVersion": "1.3",
          "certificates": [
            {
              "certificateFile": "",
              "keyFile": ""
            }
          ]
        }
      },
      "sniffing": {
        "enabled": true,
        "destOverride": [
          "http",
          "tls"
        ]
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "tag": "direct"
    },
    {
      "protocol": "blackhole",
      "tag": "block"
    }
  ],
  "policy": {
    "levels": {
      "0": {
        "handshake": 2,
        "connIdle": 120
      }
    }
  }
}`

type XrayConfig struct {
	Log struct {
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			Domain      []string `json:"domain,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			Ip          []string `json:"ip,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	Inbounds []struct {
		Listen   string `json:"listen"`
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
		Settings struct {
			Clients []struct {
				Id   string `json:"id"`
				Flow string `json:"flow"`
			} `json:"clients"`
			Decryption string `json:"decryption"`
			Fallbacks  []struct {
				Dest string `json:"dest"`
				Xver int    `json:"xver"`
			} `json:"fallbacks"`
		} `json:"settings"`
		StreamSettings struct {
			Network     string `json:"network"`
			Security    string `json:"security"`
			TlsSettings struct {
				RejectUnknownSni bool   `json:"rejectUnknownSni"`
				MinVersion       string `json:"minVersion"`
				Certificates     []struct {
					CertificateFile string `json:"certificateFile"`
					KeyFile         string `json:"keyFile"`
				} `json:"certificates"`
			} `json:"tlsSettings"`
		} `json:"streamSettings"`
		Sniffing struct {
			Enabled      bool     `json:"enabled"`
			DestOverride []string `json:"destOverride"`
		} `json:"sniffing"`
	} `json:"inbounds"`
	Outbounds []struct {
		Protocol string `json:"protocol"`
		Tag      string `json:"tag"`
	} `json:"outbounds"`
	Policy struct {
		Levels struct {
			Field1 struct {
				Handshake int `json:"handshake"`
				ConnIdle  int `json:"connIdle"`
			} `json:"0"`
		} `json:"levels"`
	} `json:"policy"`
}
