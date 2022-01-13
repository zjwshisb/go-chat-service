package util

import (
	"github.com/spf13/viper"
	"log"
	"net"
)

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln("unable to get mac address")
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

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

func GetEnv() string {
	env :=  viper.GetString("App.Env")
	if env == "" {
		env = "local"
	}
	return env
}

func GetStoragePath() string  {
	return GetWorkDir() + "/storage"
}

func GetWorkDir() string {
	workDir := viper.GetString("App.WorkDir")
	if workDir == "" {
		workDir = "./"
	}
	return workDir
}