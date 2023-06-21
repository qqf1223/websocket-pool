package nodeManager

import (
	"context"
	"encoding/json"
	"strings"
	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/pkg/discover/getcdv3"
	"websocket-pool/pkg/request"
	"websocket-pool/pkg/utils"

	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type NodeManager struct {
}

var NodeM NodeManager

func (n *NodeManager) RegisterNodeInfo() (err error) {
	registerIP, _ := utils.GetOutboundIP()
	userNode := &entity.NodeEntity{
		NodeName: global.GVA_CONFIG.System.Name,
		Ip:       registerIP,
		Port:     12345,
	}
	err = getcdv3.RegisterEtcd(userNode)
	if err != nil {
		global.GVA_LOG.Error("RegisterEtcd failed ", zap.Error(err))
		return
	}

	global.GVA_LOG.Info("RegisterEtcd ", zap.String("schema", global.GVA_CONFIG.Etcd.Schema), zap.String("addr", strings.Join(global.GVA_CONFIG.Etcd.Addr, ",")), zap.String("registerIP", registerIP), zap.String("srvName", global.GVA_CONFIG.System.Name))
	return
}

func (n *NodeManager) Distribute(data entity.MessageEntity) (err error) {
	nmsg := entity.NodeMessageEntity{
		Context: data.BizContext,
		Body:    data.Body,
		Async:   true,
	}
	nmsg.Type = "msg"
	body, err := json.Marshal(nmsg)
	if err != nil {
		global.GVA_LOG.Error("json Marshal error ", zap.Error(err))
		return
	}

	//
	srvs, err := n.GetServers(data.Context)
	if err != nil {
		global.GVA_LOG.Error("GetServers failed ", zap.Error(err))
		return
	}
	registerIP, _ := utils.GetOutboundIP()
	for _, v := range srvs {
		if v.Ip != registerIP {
			n.publisherToNodes(body, "udp", true, v)
		}
	}

	return nil
}

func (n *NodeManager) GetServers(ctx context.Context) (srvs []*entity.NodeEntity, err error) {
	return getcdv3.GetServers(ctx)
}

func (n *NodeManager) publisherToNodes(body []byte, schema string, async bool, srv *entity.NodeEntity) {
	host := srv.Ip + ":" + cast.ToString(srv.Port)
	if schema == "udp" {
		global.GVA_LOG.Info("udp发送分发命令", zap.String("body", string(body[0:])), zap.String("toIP", srv.Ip))
		request.Udp.Send(host, string(body[0:]))
	}
}
