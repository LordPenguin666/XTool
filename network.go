package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"time"
)

func (c *Config) IPAddress() *Config {
	client := resty.New()
	client.SetTimeout(time.Minute)
	client.SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	client.SetHeader("referer", "https://www.bilibili.com/")

	result := &IPZoneResp{}

	r, err := client.R().
		SetResult(result).
		Get("https://api.bilibili.com/x/web-interface/zone")

	if err != nil {
		c.logger.Error(err.Error())
		return c
	}

	if result.Code != 0 || r.StatusCode() != 200 {
		c.logger.Error("", data("response", r.String())...)
		return c
	}

	c.IP = result.Data.Addr

	return c
}

func (c *Config) DomainLookup() *Config {
	ips, err := net.LookupIP(c.Domain)
	if err != nil {
		c.logger.Error(err.Error())
	}

	var exist bool
	for _, ip := range ips {
		if ip.String() == c.IP {
			exist = true
			break
		}
	}

	if !exist {
		c.logger.Warn("please confirm your domain resolved is correct",
			data(
				"domain", c.Domain,
				"current_ip", c.IP,
				"lookup", ips)...)
		fmt.Printf("%v %v\n",
			red("[Warning]"),
			yellow("请注意, 必须将域名解析到本机上, 才能成功申请到证书!"),
		)
	} else {
		fmt.Printf("%v %v -> %v\n", green("域名解析正确"), blue(c.Domain), blue(c.IP))
	}

	return c

}

type IPZoneResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Addr        string  `json:"addr"`
		Country     string  `json:"country"`
		Province    string  `json:"province"`
		City        string  `json:"city"`
		Isp         string  `json:"isp"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		ZoneId      int     `json:"zone_id"`
		CountryCode int     `json:"country_code"`
	} `json:"data"`
}
