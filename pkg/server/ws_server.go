package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

var clientManager = client.NewClientManager()

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
		WriteBufferSize:  global.GVA_CONFIG.WS.MaxMsgLen,
		WriteBufferPool:  nil,
		CheckOrigin: func(r *http.Request) bool {
			return r.Method == http.MethodGet
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
	// 开启客户端连接管理
	go clientManager.Start()
	global.GVA_LOG.Info(fmt.Sprintf("websocket server start sucess.listen:%s", ws.addr))
	err := srv.ListenAndServe() //Start listening
	if err != nil {
		return fmt.Errorf("Ws listening err: %s" + err.Error())
	}

	return nil
}

func (ws *WsServer) WebsocketEntry(ctx *gin.Context) {
	ws.headerCheck(ctx)
	conn, err := ws.upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// handler error
		global.GVA_LOG.Error("connect error", zap.Error(err))
		return
	} else {
		newClient := client.NewClient(conn, conn.RemoteAddr().String())
		clientManager.Register <- newClient
		go ws.readMsg(ctx, newClient)
		// go ws.writeMsg()
	}

}

func (ws *WsServer) readMsg(ctx *gin.Context, c *client.Client) {
	for {
		token := ctx.GetHeader("token")
		fmt.Println(token)
		messageType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			global.GVA_LOG.Error("ws readMsg error ", zap.String("userIP", c.Addr), zap.Error(err))
			ws.delClientConn(c)
			return
		}
		if messageType == websocket.PingMessage {
			global.GVA_LOG.Info("recive ping message")
			return
		}
		if messageType == websocket.CloseMessage {
			global.GVA_LOG.Info("close message")
			ws.delClientConn(c)
			return
		}

		ws.msgParse(c, msg)
	}
}

func (ws *WsServer) writeMsg(c *client.Client, t int, msg []byte) {

	c.Conn.WriteMessage(t, msg)
}

func (ws *WsServer) delClientConn(c *client.Client) {
	err := c.Conn.Close()
	if err != nil {
		global.GVA_LOG.Error("close err", zap.Error(err))
		return
	}
	clientManager.UnRegister <- c
}

func (ws *WsServer) msgParse(c *client.Client, msg []byte) {
	global.GVA_LOG.Info("receive message", zap.String("msg", cast.ToString(msg)))
	req := &entity.Req{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		global.GVA_LOG.Error("json Unmarshal req msg error", zap.Error(err))
		return
	}
	resp := &entity.Resp{}
	result, err := json.Marshal(resp)
	if err != nil {
		global.GVA_LOG.Error("json Marshal resp msg error", zap.Error(err))
		return
	}

	ws.writeMsg(c, websocket.TextMessage, result)
}

func (ws *WsServer) headerCheck(ctx *gin.Context) {
	token, _ := ctx.GetQuery("token")
	sendID, _ := ctx.GetQuery("sendID")

	fmt.Println(token)
	fmt.Println(sendID)
}
