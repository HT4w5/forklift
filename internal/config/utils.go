package config

import (
	"strings"
	"time"

	"github.com/docker/go-units"
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

type ByteSize int64

func (b *ByteSize) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	res, err := units.FromHumanSize(s)
	if err != nil {
		return err
	}

	*b = ByteSize(res)
	return nil
}
