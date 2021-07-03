package routers

import (
	"github.com/1k-ct/vtuber-cho/handler"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/vtubers/:affiliations/:types", handler.HandlerRandItem)
	// 検索 query
	v1.GET("/search/:channel", handler.HandlerItemSearch)
	return r
}
