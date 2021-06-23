package routers

import (
	"github.com/1k-ct/twitter-dem/pkg/middelware"
	"github.com/1k-ct/vtuber-cho/backend/handler"

	// "github.com/1k-ct/vtuber-cho/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(secretKey string) *gin.Engine {
	r := gin.Default()

	// admin

	// jwtを所得
	// /admin?API_KEY=******
	r.GET("/api/v1/admin", handler.FitchJwt)

	// jwt check group --middleware--
	// |- vtuber type 登録
	r.POST("/api/v1/vtubers", middelware.TokenAuthMiddleware(secretKey), handler.RegisterVtuber)

	// r.POST("/api/v1/vtubers-file", handler.RegisterVtuberJsonFile)
	// users

	r.GET("/api/v1/vtubers/:affiliations/:types", handler.FitchRandVtuber)

	// 検索 query
	r.GET("api/v1/search", handler.SearchVtuber)

	return r
}
