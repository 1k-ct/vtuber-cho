package main

import (
	"log"

	"github.com/1k-ct/vtuber-cho/routers"
)

// docker-compose up --build -d
// docker-compose down
// docker-compose logs

func main() {
	if err := routers.Init().Run(":8000"); err != nil {
		log.Println(err)
	}
}
