package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Log     LogConfig     `json:"log"`
	Profile ProfileConfig `json:"profile"`
	Patches []PatchConfig `json:"patches"`
}

type LogConfig struct {
	Level  string       `json:"level"`
	Output string       `json:"output"`
	Syslog SyslogConfig `json:"syslog"`
}

type ProfileConfig struct {
	URL     string   `json:"url"`
	Update  Duration `json:"update"`
	Patches []string `json:"patches"`
	LogCap  struct {
		Enabled bool         `json:"enabled"`
		Output  string       `json:"output"`
		Syslog  SyslogConfig `json:"syslog"`
	}
}

type PatchConfig struct {
	Tag     string `json:"tag"`
	Content any    `json:"content"`
}

func (c *Config) Load(path string) error {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return err
	}
	return nil
}
