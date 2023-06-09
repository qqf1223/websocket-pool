package client

import (
	"context"
	"time"
	"websocket-pool/entity"
	"websocket-pool/global"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Context   entity.ContextEntity
	Conn      *websocket.Conn
	Addr      string
	Jobs      chan string
	Timestamp int64
	Key       string
}

func NewClient(ctx context.Context, ws *websocket.Conn, wsObj *entity.Req) *Client {
	context := entity.ContextEntity{
		Context:    ctx,
		AppID:      wsObj.AppID,
		PlatformID: wsObj.PlatformID,
		Token:      wsObj.Token,
		GID:        wsObj.GID,
	}
	client := &Client{
		Context:   context,
		Key:       entity.GetContextKey(context),
		Conn:      ws,
		Addr:      ws.RemoteAddr().String(),
		Jobs:      make(chan string, 100),
		Timestamp: time.Now().UnixNano(),
	}
	client.init()
	return client
}

// 初始化向客户端发送数据的缓存
func (c *Client) init() {
	go c.sandRunner()
}

func (c *Client) sandRunner() {
	for {
		select {
		case job, ok := <-c.Jobs:
			if !ok {
				return
			}
			c.safeSend(job)
		}
	}
}

func (c *Client) safeSend(msg string) {
	defer func() {
		if err := recover(); err != nil {
			global.GVA_LOG.Error("发送消息panic", zap.Any("err", err))
		}
	}()
	t := time.Now()
	err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		global.GVA_LOG.Error("向客户端发送数据异常", zap.Error(err))
	}
	global.GVA_LOG.Info("conn send end", zap.String("key", c.Key), zap.Duration("cost", time.Since(t)))
}

func (c *Client) SendTo(msg entity.MessageEntity) {
	c.Jobs <- msg.Body
}
