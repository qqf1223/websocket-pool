package entity

import "fmt"

type MessageEntity struct {
	Context ContextEntity // 公参
	Body    string
}

type MessageBody struct {
	From        string `json:"from"  validate:"required"` // 发送方
	Cmd         string `json:"cmd"  validate:"required"`
	Data        string `json:"data"`      // 消息内容
	Timestamp   int64  `json:"timestamp"` // 操作时间
	OperationID string //
}
type ContextEntity struct {
	AppID      string `json:"appId"`
	PlatformID string `json:"platformId"`
	Token      string `json:"token"`
	GID        string `json:"gid"`
}

func GetContextKey(ctx ContextEntity) string {
	key := fmt.Sprintf("%s_%s_%s", ctx.AppID, ctx.PlatformID, ctx.GID)
	return key
}

type BizMessageEntity struct {
	AppID       string `json:"appId"  validate:"required"`
	PlatformID  string `json:"platformId"  validate:"required"`
	From        string `json:"from"  validate:"required"` // 发送方
	Cmd         string `json:"cmd"  validate:"required"`
	Data        string `json:"data"`      // 消息内容
	Timestamp   int64  `json:"timestamp"` // 操作时间
	OperationID string //
	Token       string `json:"token"`
	GID         string `json:"gid"`
}
