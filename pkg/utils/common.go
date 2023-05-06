package utils

import "websocket-pool/global"

func IsDevelopment() bool {
	if global.GVA_CONFIG.System.Env == "develop" {
		return true
	}

	return false
}
