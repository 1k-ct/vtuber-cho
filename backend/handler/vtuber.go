package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/1k-ct/twitter-dem/pkg/appErrors"
	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
)

type VtuberType struct {
	gorm.Model
	VtuberID uint
	Types    string
}
type Vtuber struct {
	gorm.Model
	Name        string
	ChannelID   string
	Affiliation string
	VtuberTags  []VtuberTag `gorm:"ForeingKey:VtuberTagID"`
}
type VtuberTag struct {
	gorm.Model
	VtuberID uint
	Tag      string
}

func FetchDatabaseVtuber(affiliation, tag string) ([]Vtuber, error) {
	// database connection localhost
	// -----------------
	config, err := database.NewLocalDB("user", "password", "sample")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	db, err := config.Connect()
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	defer db.Close()
	// -------------------
	vtubers := []Vtuber{}
	if err := db.Model(&Vtuber{}).Where("affiliation = ?", affiliation).Preload("VtuberTags", "tag = ?", tag).Find(&vtubers).Error; err != nil {
		return nil, xerrors.New(err.Error())
	}
	for i, v := range vtubers {
		if len(v.VtuberTags) == 0 {
			vtubers = append(vtubers[:i], vtubers[i+1:]...)
		}
	}
	return vtubers, nil
}
func RegisterDatabaseVtuber(v *Vtuber) (*Vtuber, error) {
	// database connection localhost
	// -----------------
	config, err := database.NewLocalDB("user", "password", "sample")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	db, err := config.Connect()
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	defer db.Close()
	// -------------------
	if err := db.Create(&v).Error; err != nil {
		return nil, xerrors.New(err.Error())
	}

	return v, nil
}

// FitchRandVtuber /:affiliations/:types 条件に合ったvtuberをランダムに紹介する
func FitchRandVtuber(c *gin.Context) {
	vAffiliation := c.Param("affiliations")
	vType := c.Param("types")

	db, err := DatabaseConnection()
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	defer db.Close()

	// // db　から条件に合ったvtuberを取得
	vtubers, err := FetchDatabaseVtuber(vAffiliation, vType)
	if err != nil {
		c.JSON(400, gin.H{"status": "not found affiliation"})
		return
	}

	if len(vtubers) == 0 {
		c.JSON(400, gin.H{"status": "not found vtuber"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(vtubers))

	url := "https://www.youtube.com/channel/" + vtubers[n].ChannelID
	c.JSON(200, gin.H{"name": vtubers[n].Name, "url": url})
}
func SearchVtuber(c *gin.Context) {
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
