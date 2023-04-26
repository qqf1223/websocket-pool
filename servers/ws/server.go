package ws

// import (
// 	"net/http"
// 	"time"

// 	"websocket-pool/global"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// var userConn uint64

// type UserConn struct {
// 	*websocket.Conn
// }
// type WsServer struct {
// 	wsAddr       string
// 	wsMaxConnNum int
// 	wsUpGrader   *websocket.Upgrader
// }

// var Ws WsServer

// func (ws *WsServer) Init() {
// 	ws.wsAddr = global.GVA_CONFIG.WS.Addr
// 	ws.wsMaxConnNum = global.GVA_CONFIG.WS.MaxConnNum
// 	ws.wsUpGrader = &websocket.Upgrader{
// 		HandshakeTimeout: time.Duration(global.GVA_CONFIG.WS.Timeout) * time.Second,
// 		ReadBufferSize:   global.GVA_CONFIG.WS.MaxMsgLen,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 		EnableCompression: false,
// 	}
// }

// func (ws *WsServer) Run() {
// 	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
// 	if err != nil {
// 		panic("Ws listening err:" + err.Error())
// 	}
// }

// func (ws *WsServer) WebsocketEntry(ctx *gin.Context) {
// 	conn, err := ws.wsUpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		// handler error
// 		return
// 	} else {
// 		newConn := &UserConn{
// 			Conn: conn,
// 		}
// 		go ws.readMsg(newConn)
// 	}

// }

// func (ws *WsServer) readMsg(conn *UserConn) {
// 	messageType, msg, err := conn.ReadMessage()
// 	if err != nil {
// 		//log
// 		userConn--
// 		ws.delUserConn(conn)
// 		return

// 	}
// 	if messageType == websocket.PingMessage {
// 		// log
// 	}
// 	if messageType == websocket.CloseMessage {
// 		// log
// 		userConn--
// 		ws.delUserConn(conn)
// 		return
// 	}

// 	ws.msgParse(conn, msg)
// }

// func (ws *WsServer) delUserConn(conn *UserConn) {

// }

// func (ws *WsServer) msgParse(conn *UserConn, msg []byte)
