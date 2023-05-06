package entity

type Req struct {
	Seq  string      `json:"seq"` // 消息唯一ID
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data,omitempty"`
}

type Header struct {
	Token string `json:"token"`
}
