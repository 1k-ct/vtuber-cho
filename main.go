package main

import (
	"log"
	"os"

	"github.com/1k-ct/vtuber-cho/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// r := router.InitRouters(os.Getenv("SECRET_KEY"))
	// if err := r.Run(); err != nil {
	// 	return
	// }
	r := routers.NewRouter(os.Getenv("SECRET_KEY"))
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	// config, err := database.NewLocalDB("user", "password", "vtuber")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db, err := config.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db.AutoMigrate(Vtuber{})
	// db.AutoMigrate(VtuberType{}).AddForeignKey("vtuber_id", "vtubers(id)", "CASCADE", "CASCADE")
	// if err := db.Where("").Error; err != nil {
	// 	log.Fatal(err)
	// }

	// chID := "UCBQd84IW8OvM8H5jftHdvmw"
	// v := &Vtuber{
	// 	// Model:       gorm.Model{},
	// 	Name:        "Kurea",
	// 	ChannelID:   chID,
	// 	Affiliation: "other",
	// }
	// if err := db.Create(v).Error; err != nil {
	// 	log.Fatal(err)
	// }
	// vt := &VtuberType{
	// 	VtuberID: v.ID,
	// 	Types:    "game",
	// }
	// if err := db.Create(vt).Error; err != nil {
	// 	log.Fatal(err)
	// }
	// if err := db.Where("id = ?", chID).First(&Vtuber{}).Error; err != nil {
	// 	log.Fatal(err)
	// }
	// v := &VtuberType{}
	// // var v Vtuber
	// db.Where("types = ?", "song").Find(&v)
	// // res.RowsAffected
	// fmt.Println(v)
	// vs := []*VtuberType{}
	// if err := db.Model(&v).Where("types = ?", "game").Find(&vs).Error; err != nil {
	// 	log.Fatal(err)
	// }
	// for _, vv := range vs {
	// 	fmt.Println(vv.ID, vv.VtuberID, vv.Types)
	// }
	// fmt.Println(vt)
}

type Vtuber struct {
	gorm.Model
	Name        string
	ChannelID   string
	Affiliation string
}
type VtuberType struct {
	gorm.Model
	VtuberID uint
	Types    string
}
