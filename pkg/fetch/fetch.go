package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	headerNameUserAgent = "User-Agent"
)

func GetProfileWithUA(ctx context.Context, url string, ua string) (any, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerNameUserAgent, ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	profileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile any
	err = json.Unmarshal(profileBytes, &profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
