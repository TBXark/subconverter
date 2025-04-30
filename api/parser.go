package api

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(line string, params *SubRequestParams) (bool, any, error)
}
type ShadowSockParser struct{}

func (r *ShadowSockParser) Parse(line string, params *SubRequestParams) (bool, any, error) {
	if !strings.HasPrefix(line, "ss://") {
		return false, line, nil
	}
	if !strings.Contains(line, ".") {
		lineDecode, err := base64.URLEncoding.DecodeString(line[5:])
		if err != nil {
			return true, nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		line = string(lineDecode)
	}
	u, err := url.Parse(line)
	if err != nil {
		return true, nil, fmt.Errorf("invalid URL format: %w", err)
	}
	userInfo, err := base64.StdEncoding.DecodeString(u.User.Username())
	if err != nil {
		return true, nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	parts := strings.Split(string(userInfo), ":")
	if len(parts) != 2 {
		return true, nil, fmt.Errorf("invalid user info format")
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return true, nil, fmt.Errorf("invalid port format: %w", err)
	}
	cipher, password := parts[0], parts[1]
	name := strings.TrimSpace(u.Fragment)
	conf := ShadowSocks{
		Name:     name,
		Type:     "ss",
		Server:   u.Hostname(),
		Port:     port,
		Cipher:   cipher,
		Password: password,
		Udp:      params.UDP,
	}
	return true, conf, nil
}
