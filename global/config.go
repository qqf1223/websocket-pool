package global

type Config struct {
	WS    WS    `mapstructure:"websocket" json:"websocket" yaml:"websocket"`
	REDIS REDIS `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type WS struct {
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	Timeout    int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	MaxMsgLen  int    `mapstructure:"maxMsgLen" json:"maxMsgLen" yaml:"maxMsgLen"`
	MaxConnNum int    `mapstructure:"maxConnNum" json:"maxConnNum" yaml:"maxConnNum"`
}

type REDIS struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Pass     string `mapstructure:"pass" json:"pass" yaml:"pass"`
	Database int    `mapstructure:"database" json:"database" yaml:"database"`
}
