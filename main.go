package main

import (
	"io"
	"log"
	"os"

	"github.com/1k-ct/vtuber-cho/handler"
	"github.com/gin-gonic/gin"
)

// docker-compose up --build -d
// docker-compose down
// docker-compose logs

func main() {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/vtubers/:affiliations/:types", handler.HandlerRandItem)
	// 検索 query
	v1.GET("/search/:channel", handler.HandlerItemSearch)

	if err := r.Run(":8000"); err != nil {
		log.Println(err)
	}
}
