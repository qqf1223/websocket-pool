package server

import (
	"log"
	"net"
	"websocket-pool/global"

	"go.uber.org/zap"
)

type UDPServer struct {
	Recv chan []byte
}

func (u *UDPServer) Run() (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", global.GVA_CONFIG.Udp.Addr)
	if err != nil {
		global.GVA_LOG.Error("创建udp异常", zap.Error(err))
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		global.GVA_LOG.Error("监听udp异常", zap.Error(err))
	}
	log.Printf("[INFO] Udp server start, listen: %s\n", global.GVA_CONFIG.Udp.Addr)
	for {
		var buf [20 * 1024]byte
		len, _, err := udpConn.ReadFromUDP(buf[0:])
		if err != nil {
			global.GVA_LOG.Error("读取udp数据异常", zap.Error(err))
			continue
		}
		u.Recv <- buf[0:len]
	}
}
