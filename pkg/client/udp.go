package client

import (
	"net"
	"websocket-pool/global"

	"go.uber.org/zap"
)

type UdpClient struct {
}

func (u *UdpClient) Send(host string, body string) (err error) {
	conn, err := net.Dial("udp", host)
	if err != nil {
		global.GVA_LOG.Error("udp dial abnormal.", zap.Error(err))
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(body))
	global.GVA_LOG.Error("udp write msg abnormal.", zap.Error(err))
	return
}
