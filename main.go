package main

import (
	"log"
	"os"

	"github.com/1k-ct/vtuber-cho/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	r := routers.NewRouter(os.Getenv("SECRET_KEY"))
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
