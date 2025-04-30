package api

import (
	"net/http"

	"github.com/TBXark/subconverter/converter"
	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	engine := gin.Default()
	convert := converter.NewConverter()
	engine.GET("/sub", func(ctx *gin.Context) {
		var params converter.ConvertParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := convert.Convert(&params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.String(http.StatusOK, res)
	})
	engine.Handler().ServeHTTP(w, r)
}
