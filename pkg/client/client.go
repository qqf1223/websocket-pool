package client

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	Addr      string
	Send      chan []byte
	Timestamp int64
}

func NewClient(ws *websocket.Conn, addr string) *Client {
	return &Client{
		Conn:      ws,
		Addr:      addr,
		Send:      make(chan []byte, 100),
		Timestamp: time.Now().UnixNano(),
	}
}
