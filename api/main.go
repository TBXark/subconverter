package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Loader interface {
	Load(r *http.Request) ([]byte, error)
}
type Decoder interface {
	Decode(data []byte) ([]string, error)
}
type Render interface {
	IsValid(line string) bool
	Render(line string) (string, error)
}
type Merger interface {
	Merge(lines []string) (string, error)
}

type HTTPLoader struct{}

func (l *HTTPLoader) Load(r *http.Request) ([]byte, error) {
	query := r.URL.Query()
	targetURL := query.Get("url")
	if targetURL == "" {
		return nil, fmt.Errorf("missing url parameter")
	}
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target URL: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return data, nil
}

type Base64Decoder struct{}

func (d *Base64Decoder) Decode(data []byte) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(decoded)), "\n")
	return lines, nil
}

type ShadowSockSurgeRender struct {
}

func (r *ShadowSockSurgeRender) IsValid(line string) bool {
	return strings.HasPrefix(line, "ss://")
}

func (r *ShadowSockSurgeRender) Render(line string) (string, error) {
	u, err := url.Parse(line)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}
	userInfo, err := base64.StdEncoding.DecodeString(u.User.Username())
	if err != nil {
		return "", fmt.Errorf("failed to decode user info: %w", err)
	}
	parts := strings.Split(string(userInfo), ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid user info format")
	}
	method, password := parts[0], parts[1]
	remark := u.Fragment
	if remark == "" {
		return "", fmt.Errorf("empty remark")
	}
	return fmt.Sprintf("%s=ss, %s,%s,encrypt-method=%s,password=\"%s\"",
		remark, u.Hostname(), u.Port(), method, password), nil
}

type SurgeMerger struct {
}

func (r *SurgeMerger) Merge(lines []string) (string, error) {
	output := strings.Builder{}
	output.WriteString("[Proxy]")
	for _, line := range lines {
		output.WriteString("\n")
		output.WriteString(line)
	}
	return output.String(), nil
}

// ProxyHandler 组合各个组件处理请求
type ProxyHandler struct {
	loader  Loader
	decoder Decoder
	render  []Render
	merger  Merger
}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		loader:  &HTTPLoader{},
		decoder: &Base64Decoder{},
		render: []Render{
			&ShadowSockSurgeRender{},
		},
		merger: &SurgeMerger{},
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.loader.Load(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lines, err := h.decoder.Decode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		for _, render := range h.render {
			if render.IsValid(line) {
				if renderedLine, e := render.Render(line); e == nil {
					renderedLines = append(renderedLines, renderedLine)
				}
				break
			}
		}
	}
	output, err := h.merger.Merge(renderedLines)
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
