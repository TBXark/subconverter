package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Loader interface {
	Load(remote string) ([]string, error)
}
type Parser interface {
	Parse(line string) (bool, any, error)
}
type Generator interface {
	Generate(proxies []any) (string, error)
}

type ShadowSocks struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Server            string `json:"server"`
	Port              int    `json:"port"`
	Cipher            string `json:"cipher"`
	Password          string `json:"password"`
	Udp               bool   `json:"udp"`
	UdpOverTcp        bool   `json:"udp-over-tcp"`
	UdpOverTcpVersion int    `json:"udp-over-tcp-version"`
	IpVersion         string `json:"ip-version"`
	Plugin            string `json:"plugin"`
	PluginOpts        struct {
		Mode string `json:"mode"`
	} `json:"plugin-opts"`
	Smux struct {
		Enabled bool `json:"enabled"`
	} `json:"smux"`
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

type ShadowSockParser struct{}

func (r *ShadowSockParser) Parse(line string) (bool, any, error) {
	if !strings.HasPrefix(line, "ss://") {
		return false, line, nil
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
	conf := ShadowSocks{
		Name:              strings.TrimSpace(u.Fragment),
		Type:              "ss",
		Server:            u.Hostname(),
		Port:              port,
		Cipher:            cipher,
		Password:          password,
		Udp:               false,
		UdpOverTcp:        false,
		UdpOverTcpVersion: 0,
		IpVersion:         "",
		Plugin:            "",
		PluginOpts: struct {
			Mode string `json:"mode"`
		}{},
		Smux: struct {
			Enabled bool `json:"enabled"`
		}{},
	}
	return true, conf, nil
}

type SurgeGenerator struct{}

func (r *SurgeGenerator) Generate(lines []any) (string, error) {
	output := strings.Builder{}
	output.WriteString("[Proxy]")
	for _, line := range lines {
		output.WriteString("\n")
		switch v := line.(type) {
		case ShadowSocks:
			output.WriteString(fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=\"%s\"",
				v.Name,
				v.Server,
				v.Port,
				v.Cipher,
				v.Password,
			))
		}
	}
	return output.String(), nil
}

type ClashGenerator struct{}

func (r *ClashGenerator) Generate(lines []any) (string, error) {
	output := strings.Builder{}
	output.WriteString("proxies:\n")
	for _, line := range lines {
		bytes, err := json.Marshal(line)
		if err != nil {
			continue
		}
		output.WriteString("- ")
		output.Write(bytes)

	}
	return output.String(), nil
}

type ProxyHandler struct {
	loader    Loader
	parser    []Parser
	generator map[string]Generator
}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		loader: &HttpBase64Loader{},
		parser: []Parser{
			&ShadowSockParser{},
		},
		generator: map[string]Generator{
			"surge": &SurgeGenerator{},
			"clash": &ClashGenerator{},
		},
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	targetURL := query.Get("url")
	targetType := query.Get("target")
	if targetURL == "" || targetType == "" {
		http.Error(w, "target url or target type is empty", http.StatusBadRequest)
		return
	}
	generator, ok := h.generator[strings.ToLower(targetType)]
	if !ok {
		http.Error(w, "target type not supported", http.StatusBadRequest)
		return
	}
	lines, err := h.loader.Load(targetURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nodes := make([]any, 0, len(lines))
	for _, line := range lines {
		for _, parser := range h.parser {
			ok, node, pErr := parser.Parse(line)
			if !ok {
				continue
			}
			if pErr != nil {
				continue
			}
			nodes = append(nodes, node)
		}
	}
	output, err := generator.Generate(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(output))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := NewProxyHandler()
	handler.ServeHTTP(w, r)
}
