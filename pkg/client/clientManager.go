package client

import (
	"context"
	"sync"

	"websocket-pool/global"
	"websocket-pool/pkg/gredis"

	"go.uber.org/zap"
)

type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	ClientCount int64              // 客户端统计
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	UnRegister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
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
		cm.Clients[c] = true
		cm.ClientCount++
	}
}

// 删除客户端
func (cm *ClientManager) DelClients(client *Client) {
	cm.ClientsLock.Lock()
	defer cm.ClientsLock.Unlock()

	if _, ok := cm.Clients[client]; ok {
		delete(cm.Clients, client)
		gredis.DoWithContext(context.Background(), "INCR", "clientCount", 1)
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
		}
	}
}
