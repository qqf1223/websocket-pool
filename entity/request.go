package entity

type Req struct {
	AppID      string `json:"appId" validate:"required"` // 业务ID
	Token      string `json:"token" validate:"required"` // 用户
	RoomID     string `json:"roomId"`                    // 本次会话ID
	PlatformID string `json:"platformID,omitempty"`      //平台
}
