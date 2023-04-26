package gredis

import (
	"crypto/tls"
	"time"
	"websocket-pool/global"
)

type Config struct {
	Host        string
	Port        int
	Password    string
	Database    int
	Timeout     time.Duration
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	TlsConfig   *tls.Config
	UseTLS      bool
	SkipVerify  bool
}

func NewConfig() *Config {
	return &Config{
		Host:        global.GVA_CONFIG.REDIS.Addr,
		Port:        global.GVA_CONFIG.REDIS.Port,
		Password:    global.GVA_CONFIG.REDIS.Pass,
		Database:    global.GVA_CONFIG.REDIS.Database,
		Timeout:     0,
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: 0,
		TlsConfig:   &tls.Config{},
		UseTLS:      false,
		SkipVerify:  false,
	}
}

func Init(c *Config) {
	once.Do(func() {
		initRedisPool(c)
	})
}

func initRedisPool(c *Config) {

}
