package main

import (
	"log"

	"github.com/1k-ct/vtuber-cho/backend/routers"
	_ "github.com/go-sql-driver/mysql"
)

// docker-compose up --build -d
// docker-compose down
// docker-compose logs

func main() {
	if err := routers.Init().Run(); err != nil {
		log.Println(err)
	}
}
