package netutil

import (
	"fmt"
	"net"
)

// LocalIp 获取本机IP 和 MAC地址
func LocalIp() (ips, macs []string) {
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
