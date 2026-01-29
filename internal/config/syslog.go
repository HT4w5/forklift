package config

type SyslogConfig struct {
	Transport string `json:"transport"`
	Addr      string `json:"addr"`
}
