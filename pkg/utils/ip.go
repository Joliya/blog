/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 16:32
 * @Description:
 */

package utils

import (
	tnet "github.com/toolkits/net"
	"net"
	"strings"
	"sync"
)

var (
	once     sync.Once
	clientIP = "127.0.0.11"
)

func GetLocalIp() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.11"
		}
	})
	return clientIP
}

func GetInternalIP() string {
	inters, err := net.Interfaces()
	if IsNotNil(err) {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if IsNotNil(err) {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if IsNotNil(ipnet.IP.To4()) {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
