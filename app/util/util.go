package util

import (
	"log"
	"net"
)

func GetIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln("unable to get ip address")
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}
