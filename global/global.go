package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	GVA_CONFIG Config
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper
)

func init() {

}

const (
	CtxKeyXSendID     = "sendID"
	CtxKeyXToken      = "token"
	CtxKeyXPlatformID = "platformID"
)
