package rpc

import (
	"context"
	"websocket-pool/protobuf"
)

type WsService struct {
	protobuf.UnimplementedWsServerServer
}

func (s *WsService) QueryUsersOnline(ctx context.Context, req *protobuf.QueryUsersOnlineReq) (resp *protobuf.QueryUsersOnlineRsp, err error) {

	return
}

func (s *WsService) SendMsg(ctx context.Context, req *protobuf.SendMsgReq) (resp *protobuf.SendMsgRsp, err error) {

	return
}

func (s *WsService) SendMsgAll(ctx context.Context, req *protobuf.SendMsgAllReq) (resp *protobuf.SendMsgAllRsp, err error) {

	return
}

func (s *WsService) GetUserList(ctx context.Context, req *protobuf.GetUserListReq) (resp *protobuf.GetUserListRsp, err error) {

	return
}
