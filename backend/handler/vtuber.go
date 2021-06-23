package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/1k-ct/twitter-dem/pkg/appErrors"
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
	vAffiliation := c.Param("affiliations")
	vType := c.Param("types")

	// config, err := database.NewLocalDB("user", "password", "vtuber")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db, err := config.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db, err := DatabaseConnection()
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	defer db.Close()

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
	// config, err := database.NewLocalDB("user", "password", "vtuber")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db, err := config.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db, err := DatabaseConnection()
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	defer db.Close()

	q := c.Query("q")

	vtuber := &Vtuber{}
	if err := db.Where("channel_id = ?", q).First(&vtuber).Error; err != nil {
		c.JSON(400, gin.H{"status": "not found affiliation"})
		return
	}
	if vtuber == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	vtuberType := &VtuberType{}
	vtuberTypes := []*VtuberType{}
	if err := db.Model(&vtuberType).Where("vtuber_id = ?", vtuber.ID).
		Find(&vtuberTypes).Error; err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ErrRecordDatabase))
		return
	}

	types := []string{}
	for _, v := range vtuberTypes {
		types = append(types, v.Types)
	}
	type resp struct {
		ID          uint
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Name        string
		ChannelID   string
		Affiliation string
		Types       []string
	}
	c.JSON(200, resp{
		ID:          vtuber.ID,
		CreatedAt:   vtuber.CreatedAt,
		UpdatedAt:   vtuber.UpdatedAt,
		Name:        vtuber.Name,
		ChannelID:   vtuber.ChannelID,
		Affiliation: vtuber.Affiliation,
		Types:       types,
	})
}
