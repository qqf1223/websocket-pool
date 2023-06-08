package client

import (
	"context"
	"sync"
	"time"

	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/pkg/gredis"

	"go.uber.org/zap"
)

type ClientManager struct {
	ClientMap   sync.Map
	Clients     map[*Client]bool     // 全部的连接
	ClientsLock sync.RWMutex         // 读写锁
	ClientCount int64                // 本地客户端统计
	Users       map[string][]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex         // 读写锁
	Register    chan *Client         // 连接连接处理
	UnRegister  chan *Client         // 断开连接处理程序
	Broadcast   chan []byte          // 广播 向全部成员发送数据
	SendC       chan entity.MessageEntity
	ReceiveC    chan entity.MessageEntity
}

var Cm ClientManager

func NewClientManager() (cm *ClientManager) {
	cm = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string][]*Client),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

func (cm *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	cm.ClientsLock.RLock()
	defer cm.ClientsLock.RUnlock()
	for k, v := range cm.Clients {
		clients[k] = v
	}
	return
}

// 新增客户端连接
func (cm *ClientManager) AddClients(c *Client) {
	cm.ClientsLock.Lock()
	defer cm.ClientsLock.Unlock()
	if _, ok := cm.Clients[c]; !ok {
		cm.ClientMap.Store(entity.GetContextKey(c.Context), c)
		cm.Clients[c] = true
		cm.ClientCount++
		gredis.DoWithContext(context.Background(), "INCR", "clientCount")
	}
}

// 删除客户端
func (cm *ClientManager) DelClients(c *Client) {
	cm.ClientsLock.Lock()
	defer cm.ClientsLock.Unlock()
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	if _, ok := cm.Clients[c]; ok {
		cm.ClientMap.Delete(entity.GetContextKey(c.Context))
		delete(cm.Clients, c)
		cm.ClientCount--
		gredis.DoWithContext(context.Background(), "DECR", "clientCount")
	}
	// ToDo 清除用户连接

	// todo清除redis数据

}

func (cm *ClientManager) EventRegister(c *Client) {
	cm.AddClients(c)
	global.GVA_LOG.Info("client connect", zap.String("userIp", c.Addr), zap.Int64("ClientCount", cm.ClientCount))
}

func (cm *ClientManager) EventUnRegister(c *Client) {
	cm.DelClients(c)
	global.GVA_LOG.Info("client disconnect", zap.String("userIp", c.Addr), zap.Int64("ClientCount", cm.ClientCount))
}

func (cm *ClientManager) Start() {
	for {
		select {
		case conn := <-cm.Register:
			cm.EventRegister(conn)
		case conn := <-cm.UnRegister:
			cm.EventUnRegister(conn)
		case msg := <-cm.SendC:
			cm.safeSend(msg)
		}

	}
}

func (cm *ClientManager) safeSend(m entity.MessageEntity) {

	// 获取client
	key := entity.GetContextKey(m.Context)
	// 若本地找到对应客户端，则发送，否则寻找客户端连接进行分发
	if conn, ok := cm.ClientMap.Load(key); ok {
		conn.(*Client).SendTo(m)
	} else {

	}
}

func (cm *ClientManager) SendTo(ctx context.Context, m *entity.BizMessageEntity) {
	// 构造消息
	to := entity.MessageEntity{
		Context: entity.ContextEntity{
			AppID:      m.AppID,
			PlatformID: m.PlatformID,
			Token:      m.Token,
			GID:        m.GID,
		},
		Body: m.Data,
	}
	t := time.Now()
	cm.SendC <- to
	global.GVA_LOG.Info("server send to client end", zap.Any("context", entity.GetContextKey(to.Context)), zap.Duration("cost", time.Since(t)))
}
