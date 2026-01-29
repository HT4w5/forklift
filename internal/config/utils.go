package config

import (
	"strings"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	du, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(du)
	return nil
}
