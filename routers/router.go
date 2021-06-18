package routers

import (
	"net/http"

	"github.com/1k-ct/twitter-dem/pkg/middelware"
	"github.com/1k-ct/vtuber-cho/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(secretKey string) *gin.Engine {
	r := gin.Default()

	// admin
	// TODO api-key と　password をハッシュ化して保存
	// TODO ハッシュ化したものと比べる
	// TODO database を変更
	// jwtを所得
	// GET api/v1/admin
	r.GET("/api/v1/admin", handler.FitchJwt)

	// jwt check group --middleware--
	// |- vtuber type 登録
	// |- post api/v1/vtubers
	r.POST("/api/v1/vtubers", middelware.TokenAuthMiddleware(secretKey), handler.RegisterVtuber)

	// users

	// api/v1/vtubers/affiliations/types
	r.GET("/api/v1/vtubers/:affiliations/:types", handler.FitchRandVtuber)

	// 検索 query
	// post api/v1/search&q=

	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//
