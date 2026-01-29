package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Log     LogConfig     `json:"log"`
	Profile ProfileConfig `json:"profile"`
	Exec    ExecConfig    `json:"exec"`
	Patches []PatchConfig `json:"patches"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type ExecConfig struct {
	Path   string `json:"path"`    // Path to sing-box binary. Defaults to "sing-box"
	LogFwd bool   `json:"log_fwd"` // Forward log to stdout
}

type ProfileConfig struct {
	URL     string   `json:"url"`
	Update  string   `json:"update"`
	UA      string   `json:"ua"`
	Patches []string `json:"patches"`
}

type PatchConfig struct {
	Tag     string `json:"tag"`
	Content any    `json:"content"`
}

func Default() *Config {
	return &Config{
		Exec: ExecConfig{
			Path: "sing-box",
		},
	}
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
