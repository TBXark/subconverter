package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Loader interface {
	Load(remote string) ([]string, error)
}

type HttpBase64Loader struct{}

func (l *HttpBase64Loader) Load(remote string) ([]string, error) {
	resp, err := http.Get(remote)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target URL: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(decoded)), "\n")
	return lines, nil
}
