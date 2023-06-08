package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no ip")
}

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "192.168.2.3:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	fmt.Println("当前ip:", strings.Split(conn.LocalAddr().String(), ":")[0])
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}
