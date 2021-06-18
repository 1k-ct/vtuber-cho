package handler

import (
	"log"
	"math/rand"
	"time"

	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

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

// FitchRandVtuber /:affiliations/:types 条件に合ったvtuberをランダムに紹介する
func FitchRandVtuber(c *gin.Context) {
	// 	/:Affiliations/:types 取得
	vAffiliation := c.Param("affiliations")
	vType := c.Param("types")
	// vAffiliation := "other"
	// vType := "song"
	config, err := database.NewLocalDB("user", "password", "vtuber")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// db　から条件に合ったvtuberを取得
	vtuber := &Vtuber{}
	vtubers := []*Vtuber{}
	if err := db.Model(&vtuber).Where("affiliation = ?", vAffiliation).Find(&vtubers).Error; err != nil {
		c.JSON(400, gin.H{"status": "not found affiliation"})
		return
	}
	vtuberIDs := []uint{}
	for _, v := range vtubers {
		vtuberIDs = append(vtuberIDs, v.ID)
	}

	vtuberType := &VtuberType{}
	vtuberTypes := []*VtuberType{}
	if err := db.Model(&vtuberType).Where("types = ?", vType).Find(&vtuberTypes).Error; err != nil {
		c.JSON(400, gin.H{"status": "not found affiliation"})
		return
	}

	res := []uint{}
	for _, vid := range vtuberIDs {
		for _, v := range vtuberTypes {
			if vid == v.VtuberID {
				res = append(res, vid)
			}
		}
	}
	// fmt.Println(res)
	// 出力する
	if len(res) == 0 {
		c.JSON(400, gin.H{"status": "not found vtuber"})
		return
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(res))

	url := "https://www.youtube.com/channel/" + vtubers[res[n]-1].ChannelID
	c.JSON(200, gin.H{"name": vtubers[res[n]-1].Name, "url": url})
}
func SearchVtuber(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
