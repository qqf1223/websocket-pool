package gredis

import (
	"context"
	"crypto/tls"
	"sync"
	"time"
	"websocket-pool/global"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cast"
)

var (
	once        sync.Once
	redisClient *redis.Pool
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
		Host:        global.GVA_CONFIG.REDIS.Host,
		Port:        global.GVA_CONFIG.REDIS.Port,
		Password:    global.GVA_CONFIG.REDIS.Pass,
		Database:    global.GVA_CONFIG.REDIS.Database,
		Timeout:     time.Duration(global.GVA_CONFIG.REDIS.Timeout) * time.Second,
		MaxIdle:     global.GVA_CONFIG.REDIS.MaxIdle,
		MaxActive:   global.GVA_CONFIG.REDIS.MaxActive,
		IdleTimeout: time.Duration(global.GVA_CONFIG.REDIS.IdleTimeout) * time.Second,
		TlsConfig:   &tls.Config{},
		UseTLS:      global.GVA_CONFIG.REDIS.UseTLS,
		SkipVerify:  global.GVA_CONFIG.REDIS.SkipVerify,
	}
}

func Init(c *Config) {
	once.Do(func() {
		initRedisPool(c)
	})
}

func initRedisPool(c *Config) {
	redisHost := c.Host + ":" + cast.ToString(c.Port)
	redisPass := c.Password
	redisDB := c.Database

	timeout := time.Duration(c.Timeout)

	redisClient = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial(
				"tcp",
				redisHost,
				redis.DialPassword(redisPass),
				redis.DialDatabase(redisDB),
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout),
				redis.DialTLSConfig(c.TlsConfig),
			)
			return
		},
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: c.IdleTimeout,
		// Wait:        true,
	}
}

func getConn() (redis.Conn, error) {
	conn := redisClient.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil
}

func DoWithContext(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn, err := getConn()
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	return conn.Do(commandName, args...)
}
