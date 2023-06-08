package client

import (
	"context"
	"time"
	"websocket-pool/protobuf"

	"google.golang.org/grpc"
)

type WsGrpcClient struct {
	conn protobuf.WsServerClient
}

var (
	instance *WsGrpcClient
)

func client() *WsGrpcClient {
	if instance != nil {
		return instance
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	cc, err := grpc.DialContext(ctx, "ws")
	if err != nil {
		return nil
	}
	instance = &WsGrpcClient{
		conn: protobuf.NewWsServerClient(cc),
	}

	return instance
}

// 私聊
func (*WsGrpcClient) SendMsg(ctx context.Context) (err error) {
	var req = &protobuf.SendMsgReq{
		Seq:     "",
		AppId:   0,
		UserId:  "",
		Cms:     "",
		Type:    "",
		Msg:     "",
		IsLocal: false,
	}
	_, err = client().conn.SendMsg(ctx, req)
	return
}

// 群聊
func (*WsGrpcClient) SendMsgAll(ctx context.Context) (err error) {

	var req = &protobuf.SendMsgAllReq{
		Seq:    "",
		AppId:  0,
		UserId: "",
		Cms:    "",
		Type:   "",
		Msg:    "",
	}

	_, err = client().conn.SendMsgAll(ctx, req)
	return
}

// 定向发送一对多
func (*WsGrpcClient) SendMsgTo(ctx context.Context) (err error) {

	return
}

// 获取频道内用户列表
func (*WsGrpcClient) GetUserList(ctx context.Context) (resp *protobuf.GetUserListRsp, err error) {

	var req = &protobuf.GetUserListReq{
		AppId: 0,
	}
	resp, err = client().conn.GetUserList(ctx, req)

	return
}

// 查询在线用户
func (*WsGrpcClient) QueryUsersOnline(ctx context.Context) (resp *protobuf.QueryUsersOnlineRsp, err error) {

	var req = &protobuf.QueryUsersOnlineReq{
		AppId:  0,
		UserId: "",
	}
	resp, err = client().conn.QueryUsersOnline(ctx, req)
	return
}
