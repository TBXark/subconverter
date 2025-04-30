package converter

import (
	"errors"
	"strings"
)

type Converter struct {
	loader    Loader
	parser    []Parser
	generator map[string]Generator
}

func NewConverter() *Converter {
	return &Converter{
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

func (c *Converter) Convert(params *ConvertParams) (string, error) {
	if params.URL == "" || params.Target == "" {
		return "", errors.New("url and target are required")
	}
	generator, supported := c.generator[strings.ToLower(params.Target)]
	if !supported {
		return "", errors.New("unsupported target")
	}
	lines, err := c.loader.Load(params.URL)
	if err != nil {
		return "", err
	}
	nodes := make([]any, 0, len(lines))
	for _, line := range lines {
		for _, parser := range c.parser {
			ok, node, pErr := parser.Parse(line, params)
			if !ok {
				continue
			}
			if pErr != nil {
				continue
			}
			nodes = append(nodes, node)
		}
	}
	output, err := generator.Generate(nodes, params)
	if err != nil {
		return "", err
	}
	return output, nil
}
