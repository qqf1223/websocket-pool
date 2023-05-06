package rpc

import "websocket-pool/protobuf"

type WsService struct {
	protobuf.UnimplementedWsServerServer
}

func (s *WsService) QueryUsersOnline(req *protobuf.QueryUsersOnlineReq) (resp protobuf.QueryUsersOnlineRsp, err error) {

	return
}

func (s *WsService) SendMsg(req *protobuf.SendMsgReq) (resp protobuf.SendMsgRsp, err error) {

	return
}

func (s *WsService) SendMsgAll(req *protobuf.SendMsgAllReq) (resp protobuf.SendMsgAllRsp, err error) {

	return
}

func (s *WsService) GetUserList(req *protobuf.GetUserListReq) (resp protobuf.GetUserListRsp, err error) {

	return
}
