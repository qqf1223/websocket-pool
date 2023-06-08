package utils

import (
	"math/rand"
	"strconv"
	"time"
	"websocket-pool/global"
)

func IsDevelopment() bool {
	return global.GVA_CONFIG.System.Env == "develop"
}

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
