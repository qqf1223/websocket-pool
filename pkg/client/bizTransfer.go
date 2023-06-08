package client

import (
	"websocket-pool/entity"
	"websocket-pool/global"
)

type BizTransfer struct{}

var BizT BizTransfer

func (b *BizTransfer) Init(rcvC chan entity.MessageEntity) {
	for i := 1; i <= global.GVA_CONFIG.System.TransferPoolSize; i++ {
		go b.toTransferBiz(rcvC)
	}
}

func (b *BizTransfer) toTransferBiz(rcvC chan entity.MessageEntity) {
	for {
		select {
		case msg, ok := <-rcvC:
			if !ok {
				// 管道关闭
				break
			}
			b.transferByRule(msg)
		}
	}
}

// 业务转发，根据规则转发到对应的服务中
func (b *BizTransfer) transferByRule(msg entity.MessageEntity) {

	go b.sendToBiz()
}

func (b *BizTransfer) sendToBiz() {

}
