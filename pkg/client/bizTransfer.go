package client

import (
	"encoding/json"
	"time"
	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/internal/ruleManager"
	"websocket-pool/pkg/request"

	"go.uber.org/zap"
)

type BizTransfer struct{}

var BizT BizTransfer

func (b *BizTransfer) Init() {
	for i := 1; i <= global.GVA_CONFIG.System.TransferPoolSize; i++ {
		go b.toTransferBiz()
	}
}

func (b *BizTransfer) toTransferBiz() {
	for {
		select {
		case msg, ok := <-Cm.ReceiveC:
			if !ok {
				// 管道关闭
				break
			}
			b.transferByRule(msg)
		default:
			continue
		}
	}
}

// 业务转发，根据规则转发到对应的服务中
func (b *BizTransfer) transferByRule(msg entity.MessageEntity) {
	// 获取服务端地址
	uri, err := ruleManager.RuleM.CheckRule()
	if err != nil {
		global.GVA_LOG.Error("transfer Rule failed.", zap.Error(err))
		return
	}
	go b.sendToBiz(msg, uri)
}

func (b *BizTransfer) sendToBiz(trm entity.MessageEntity, uri string) {
	t := time.Now()
	var err error
	defer func() {
		if err == nil {
			global.GVA_LOG.Info("send to server end", zap.Any("context", entity.GetContextKey(trm.BizContext)), zap.Duration("cost", time.Since(t)))
		}
	}()
	header := make(map[string]string)
	post, _ := json.Marshal(trm)
	body, err := request.PostWithContext(trm.Context, uri, post, global.GVA_CONFIG.Http.Timeout, header)
	if err != nil {
		global.GVA_LOG.Error("http request failed.", zap.Error(err), zap.Any("context", entity.GetContextKey(trm.BizContext)))
		return
	}

	result := &entity.BizMessageEntity{}
	err = json.Unmarshal(body, result)
	if err != nil {
		global.GVA_LOG.Error("server response Unmarshal error.", zap.Error(err), zap.Any("context", entity.GetContextKey(trm.BizContext)))
		return
	}
	// 如果设置的是同步，获取到内容开始分发，否则等待服务端调用
	if global.GVA_CONFIG.System.SrvSync {
		Cm.SendTo(trm.Context, result)
	}
}
