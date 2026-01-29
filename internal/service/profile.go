package service

import (
	"context"
	"time"

	"github.com/HT4w5/song/pkg/fetch"
	"github.com/HT4w5/song/pkg/patch"
)

const (
	timeout = time.Second * 10
)

func (svc *Service) compileProfile() (any, error) {
	fetchCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	profile, err := fetch.GetProfileWithUA(
		fetchCtx,
		svc.cfg.Profile.URL,
		svc.cfg.Profile.UA,
	)

	if err != nil {
		return nil, err
	}

	// Apply patches
	for _, v := range svc.patches {
		profile = patch.Patch(profile, v)
	}

	return profile, nil
}
