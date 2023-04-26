package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"websocket-pool/global"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

var userConn uint64

type UserConn struct {
	*websocket.Conn
}
type WsServer struct {
	handler    http.Handler
	addr       string
	port       int
	maxConnNum int
	upGrader   *websocket.Upgrader
}

var Ws WsServer

func (ws *WsServer) Init(engine *gin.Engine) *WsServer {
	ws.handler = engine
	ws.addr = global.GVA_CONFIG.WS.Addr + ":" + cast.ToString(global.GVA_CONFIG.WS.Port)
	ws.port = global.GVA_CONFIG.WS.Port
	ws.maxConnNum = global.GVA_CONFIG.WS.MaxConnNum
	ws.upGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(global.GVA_CONFIG.WS.Timeout) * time.Second,
		ReadBufferSize:   global.GVA_CONFIG.WS.MaxMsgLen,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	return ws
}

func (ws *WsServer) Run() error {
	srv := &http.Server{
		Addr:    ws.addr,
		Handler: ws.handler,
	}
	err := srv.ListenAndServe() //Start listening
	if err != nil {
		return fmt.Errorf("Ws listening err: %s" + err.Error())
	}
	log.Printf("[INFO] ws server start, listen: %s\n", ws.addr)
	return nil
}

func (ws *WsServer) WebsocketEntry(ctx *gin.Context) {
	conn, err := ws.upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// handler error
		return
	} else {
		newConn := &UserConn{
			Conn: conn,
		}
		go ws.readMsg(newConn)
	}

}

func (ws *WsServer) readMsg(conn *UserConn) {
	messageType, msg, err := conn.ReadMessage()
	if err != nil {
		//log
		userConn--
		ws.delUserConn(conn)
		return

	}
	if messageType == websocket.PingMessage {
		// log
	}
	if messageType == websocket.CloseMessage {
		// log
		userConn--
		ws.delUserConn(conn)
		return
	}

	ws.msgParse(conn, msg)
}

func (ws *WsServer) delUserConn(conn *UserConn) {

}

func (ws *WsServer) msgParse(conn *UserConn, msg []byte) {

}
