package routers

import (
	"github.com/1k-ct/twitter-dem/pkg/middelware"
	"github.com/1k-ct/vtuber-cho/backend/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(secretKey string) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("v1")
	// admin

	// jwtを所得
	// /admin?API_KEY=******
	v1.GET("/admin", handler.FitchJwt)

	// jwt check group --middleware--
	// |- vtuber type 登録
	v1.POST("/vtubers", middelware.TokenAuthMiddleware(secretKey), handler.RegisterVtuber)

	// r.POST("/api/v1/vtubers-file", handler.RegisterVtuberJsonFile)
	// users

	v1.GET("/vtubers/:affiliations/:types", handler.FitchRandVtuber)

	// 検索 query
	v1.GET("/search", handler.SearchVtuber)

	return r
}
