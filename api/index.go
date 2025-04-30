package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Application struct {
	loader    Loader
	parser    []Parser
	generator map[string]Generator
}

func NewApplication() *Application {
	return &Application{
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

func (h *Application) sub(ctx *gin.Context) {
	var params SubRequestParams
	err := ctx.BindQuery(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.URL == "" || params.Target == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "url and target are required"})
		return
	}
	generator, supported := h.generator[strings.ToLower(params.Target)]
	if !supported {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported target"})
		return
	}
	lines, err := h.loader.Load(params.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nodes := make([]any, 0, len(lines))
	for _, line := range lines {
		for _, parser := range h.parser {
			ok, node, pErr := parser.Parse(line, &params)
			if !ok {
				continue
			}
			if pErr != nil {
				continue
			}
			nodes = append(nodes, node)
		}
	}
	output, err := generator.Generate(nodes, &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, output)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	engine := gin.Default()
	app := NewApplication()
	engine.GET("/sub", app.sub)
	engine.Handler().ServeHTTP(w, r)
}
