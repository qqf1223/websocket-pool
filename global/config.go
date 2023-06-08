package global

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	WS     WS     `mapstructure:"websocket" json:"websocket" yaml:"websocket"`
	Http   Http   `mapstructure:"http" json:"http" yaml:"http"`
	Rpc    Rpc    `mapstructure:"rpc" json:"rpc" yaml:"rpc"`
	REDIS  REDIS  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Udp    Udp    `mapstructure:"udp" json:"udp" yaml:"udp"`
	Etcd   Etcd   `mapstructure:"etcd" json:"etcd" yaml:"etcd"`
}
type System struct {
	Env              string `mapstructure:"env" json:"env" yaml:"env"`                                        // 环境值
	TransferPoolSize int    `mapstructure:"transferPoolSize" json:"transferPoolSize" yaml:"transferPoolSize"` // 业务转发池控制
}

type Http struct {
	Addr    string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Timeout int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
type WS struct {
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Timeout    int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	MaxMsgLen  int    `mapstructure:"maxMsgLen" json:"maxMsgLen" yaml:"maxMsgLen"`
	MaxConnNum int    `mapstructure:"maxConnNum" json:"maxConnNum" yaml:"maxConnNum"`
}
type Rpc struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type REDIS struct {
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`
	Pass        string `mapstructure:"password" json:"password" yaml:"password"`
	Database    int    `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdle     int    `mapstructure:"maxIdle" json:"maxIdle" yaml:"maxIdle"`
	MaxActive   int    `mapstructure:"maxActive" json:"maxActive" yaml:"maxActive"`
	IdleTimeout int    `mapstructure:"idleTimeout" json:"idleTimeout" yaml:"idleTimeout"`
	Timeout     int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	UseTLS      bool   `mapstructure:"useTls" json:"useTls" yaml:"useTls"`
	SkipVerify  bool   `mapstructure:"skipVerify" json:"skipVerify" yaml:"skipVerify"`
}

type Udp struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
}

// TransportLevel 根据字符串转化为 zapcore.Level
func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

type Etcd struct {
	Schema   string   `mapstructure:"schema" json:"schema" yaml:"schema"`
	Addr     []string `mapstructure:"addr" json:"addr" yaml:"addr"`
	UserName string   `mapstructure:"userName" json:"userName" yaml:"userName"`
	Password string   `mapstructure:"password" json:"password" yaml:"password"`
	Secret   string   `mapstructure:"secret" json:"secret" yaml:"secret"`
}
