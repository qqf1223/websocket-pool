package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/pkg/client"
	"websocket-pool/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// var clientManager = client.NewClientManager()

type WsServer struct {
	Handler    http.Handler
	Addr       string
	MaxConnNum int
	UpGrader   *websocket.Upgrader
}

var Ws WsServer

func (ws *WsServer) Init(engine *gin.Engine) *WsServer {
	ws.Handler = engine
	ws.Addr = global.GVA_CONFIG.WS.Addr
	ws.MaxConnNum = global.GVA_CONFIG.WS.MaxConnNum
	ws.UpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(global.GVA_CONFIG.WS.Timeout) * time.Second,
		ReadBufferSize:   global.GVA_CONFIG.WS.MaxMsgLen,
		WriteBufferSize:  global.GVA_CONFIG.WS.MaxMsgLen,
		WriteBufferPool:  nil,
		CheckOrigin: func(r *http.Request) bool {
			return r.Method == http.MethodGet
		},
		EnableCompression: false,
	}
	client.Cm = client.Cm.NewClientManager()
	// 开启客户端连接管理
	go client.Cm.Start()
	return ws
}

func (ws *WsServer) Run() error {
	srv := &http.Server{
		Addr:    ws.Addr,
		Handler: ws.Handler,
	}

	log.Printf("[INFO] websocket server start, listen: %s\n", ws.Addr)
	err := srv.ListenAndServe() //Start listening
	if err != nil {
		return fmt.Errorf("Ws listening err: %s" + err.Error())
	}

	return nil
}

func (ws *WsServer) WebsocketEntry(ctx *gin.Context) {
	isPass, wsObj, err := ws.ValidityCheck(ctx)
	if err != nil {
		global.GVA_LOG.Error("authentication error", zap.Error(err))
		return
	}
	if isPass {
		conn, e := ws.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if e != nil {
			// handler error
			global.GVA_LOG.Error("connect error", zap.Error(err))
			return
		} else {
			newClient := client.NewClient(ctx, conn, wsObj)
			client.Cm.Register <- newClient
			go ws.readMsg(ctx, newClient)
		}
	} else {
		global.GVA_LOG.Error("authentication failed", zap.Error(err))
	}
}

func (ws *WsServer) readMsg(ctx *gin.Context, c *client.Client) {
	for {
		messageType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			global.GVA_LOG.Error("ws readMsg error ", zap.String("userIP", c.Addr), zap.Error(err))
			client.Cm.UnRegister <- c
			return
		}
		if messageType == websocket.PingMessage {
			global.GVA_LOG.Info("recive ping message")
			return
		}
		if messageType == websocket.CloseMessage {
			global.GVA_LOG.Info("close message")
			client.Cm.UnRegister <- c
			return
		}

		ws.msgParse(c, msg)
	}
}

func (ws *WsServer) msgParse(c *client.Client, msg []byte) {
	global.GVA_LOG.Info("receive message", zap.String("msg", cast.ToString(msg)))
	body := &entity.MessageBody{}
	err := json.Unmarshal(msg, body)
	if err != nil {
		global.GVA_LOG.Error("json Unmarshal req msg error", zap.Error(err))
		return
	}

	// 判定是客户端-客户端 or 客户端-服务端（需要服务端进行处理或需要广播），服务端-客户端（点对点/广播）
	if body.Cmd == consts.Signalling_HeartBeat {
		// 直接回复客户端指令
		resp := fmt.Sprintf(consts.Signalling_HEARTBEAT_Resp, c.Context.PlatformID, body.Timestamp, body.Timestamp, body.Cmd)
		c.Jobs <- resp
	} else if c.Context.PlatformID != consts.Platform_Server {
		// 需要服务端处理转发
		client.Cm.ReceiveC <- entity.MessageEntity{
			BizContext: c.Context,
			Body:       string(msg),
		}
	} else if c.Context.PlatformID == consts.Platform_Server {
		// 接收服务端信令，需要寻找客户端进行处理
		client.Cm.SendC <- entity.MessageEntity{
			BizContext: entity.ContextEntity{},
			Body:       string(msg),
		}
	}
}

func (ws *WsServer) ValidityCheck(ctx *gin.Context) (isPass bool, wsObj *entity.Req, err error) {
	appId, _ := ctx.GetQuery("appId")
	token, _ := ctx.GetQuery("token")
	gid, _ := ctx.GetQuery("gid")
	platformID, _ := ctx.GetQuery("platformID")
	if appId == "" || token == "" || platformID == "" {
		return
	}

	wsObj = &entity.Req{
		AppID:      appId,
		Token:      token,
		GID:        gid,
		PlatformID: platformID,
	}
	// TODO: 通过第三方校验token
	// 暂时使用本地进行校验
	isPass = true
	return
}
