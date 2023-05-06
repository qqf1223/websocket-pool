package entity

type Message struct {
	From       int    // 发送方
	PlatformID int    // 平台
	Timestamp  int64  // 操作时间
	Content    string // 消息内容
}
