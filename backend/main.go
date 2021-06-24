package main

import (
	"log"
	"os"

	"github.com/1k-ct/vtuber-cho/backend/handler"
	"github.com/1k-ct/vtuber-cho/backend/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// docker-compose up --build -d
// docker-compose down
// docker-compose logs

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	if os.Getenv("SECRET_KEY") == "" {
		log.Fatal("not found SECRET_KEY")
	}

	r := routers.NewRouter(os.Getenv("SECRET_KEY"))
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
func createDatabases() {
	db, err := handler.DatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(handler.Vtuber{}).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(handler.VtuberType{}).
		AddForeignKey("vtuber_id", "vtubers(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.Fatal(err)

	}
}
