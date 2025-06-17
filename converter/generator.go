package converter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Generator interface {
	Generate(proxies []any, params *ConvertParams) (string, error)
}

type SurgeGenerator struct{}

func (r *SurgeGenerator) Generate(lines []any, params *ConvertParams) (string, error) {
	output := strings.Builder{}
	output.WriteString("[Proxy]")
	for _, line := range lines {
		output.WriteString("\n")
		switch v := line.(type) {
		case ShadowSocks:
			name := v.Name
			if params.AppendType {
				name = fmt.Sprintf("[ShadowSocks] %s", name)
			}
			output.WriteString(fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=\"%s\", tfo=%s, udp-relay=%s",
				name,
				v.Server,
				v.Port,
				v.Cipher,
				v.Password,
				strconv.FormatBool(params.TFO),
				strconv.FormatBool(params.UDP),
			))
		}
	}
	return output.String(), nil
}

type ClashGenerator struct{}

func (r *ClashGenerator) Generate(lines []any, params *ConvertParams) (string, error) {
	output := strings.Builder{}
	output.WriteString("proxies:\n")
	for _, line := range lines {
		bytes, err := json.Marshal(line)
		if err != nil {
			continue
		}
		output.WriteString("  - ")
		output.Write(bytes)
		output.WriteString("\n")
	}
	return output.String(), nil
}
