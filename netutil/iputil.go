package netutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

// GetLocalIp 获取本机IP 和 MAC地址
func GetLocalIp() (ips, macs []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		address, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range address {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ips = append(ips, ip.String())
			macs = append(macs, iface.HardwareAddr.String())
		}
	}
	return
}

// GetRemoteIp 快速获取本机的外网IP
func GetRemoteIp() []string {
	b, err := Get("http://httpbin.org/ip", nil)
	if err != nil {
		return []string{}
	}
	var ipMap map[string]string
	err = json.Unmarshal(b, &ipMap)
	origin := ipMap["origin"]
	if origin == "" {
		return []string{}
	}
	return strings.Split(origin, ",")
}

// GetIpToLocation 根据IP查询地理位置，默认查询本机。https://ip-api.com/
func GetIpToLocation(ip string) (*location, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=66846719&lang=zh-CN", ip)
	b, err := Get(url, nil)
	if err != nil {
		return nil, err
	}
	l := &location{}
	err = json.Unmarshal(b, l)
	if err != nil {
		return nil, err
	}
	if l.Status != "success" {
		return nil, errors.New(l.Message)
	}
	return l, nil
}

type location struct {
	Status        string  `json:"status"` // success or fail
	Message       string  `json:"message"`
	Country       string  `json:"country"`       // 国家
	CountryCode   string  `json:"countryCode"`   // 国家代码
	Continent     string  `json:"continent"`     // 洲
	ContinentCode string  `json:"continentCode"` // 洲代码
	Region        string  `json:"region"`        // 省
	RegionName    string  `json:"regionName"`    // 省代码
	City          string  `json:"city"`          // 市
	District      string  `json:"district"`      // 区
	Zip           string  `json:"zip"`           // 邮编
	Lat           float64 `json:"lat"`           // 纬度
	Lon           float64 `json:"lon"`           // 经度
	Timezone      string  `json:"timezone"`      // 时区
	Offset        int     `json:"offset"`        // Timezone UTC DST offset in seconds
	Currency      string  `json:"currency"`      // 货币
	Isp           string  `json:"isp"`           // ISP
	Org           string  `json:"org"`           // 组织
	As            string  `json:"as"`
	AsName        string  `json:"asname"`
	Reverse       string  `json:"reverse"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Query         string  `json:"query"` // 查询的IP
}
