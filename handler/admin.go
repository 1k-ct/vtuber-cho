package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/1k-ct/twitter-dem/pkg/appErrors"
	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/1k-ct/twitter-dem/pkg/jwtToken"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// type LoginInfo struct {
// 	APIKey string
// }

// type ResponseJson struct {
// 	Status  string
// 	Message string
// 	Data    interface{}
// }

func FitchJwt(c *gin.Context) {
	apiKey := c.Request.Header.Get("api-key")
	if err := godotenv.Load(); err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}

	if os.Getenv("API_KEY") != apiKey {
		c.JSON(401, gin.H{"status": "Invalid key"})
		return
	}

	ts, err := jwtToken.CreateToken(os.Getenv("USER_ID"), os.Getenv("USER_NAME"), os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_KEY"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, appErrors.ErrNotCreateToken)
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(200, tokens)
}

type RequestVtuber struct {
	Name        string   `json:"name"`
	ChannelID   string   `json:"channel_id"`
	Affiliation string   `json:"affiliation"`
	Types       []string `json:"types"`
}

func RegisterVtuber(c *gin.Context) {
	// ---------------------------
	config, err := database.NewLocalDB("user", "password", "vtuber")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// -----------------------------

	// type RequestVtuber struct {
	// 	Name        string   `json:"name"`
	// 	ChannelID   string   `json:"channel_id"`
	// 	Affiliation string   `json:"affiliation"`
	// 	Types       []string `json:"types"`
	// }
	requestVtuber := &RequestVtuber{}

	if err := c.BindJSON(&requestVtuber); err != nil {
		c.JSON(400, appErrors.ErrMeatdataMsg(err, appErrors.ErrorJSON))
		return
	}
	vtuber := &Vtuber{
		Name:        requestVtuber.Name,
		ChannelID:   requestVtuber.ChannelID,
		Affiliation: requestVtuber.Affiliation,
	}
	vtubers := []*Vtuber{}
	res := db.Where("channel_id = ?", requestVtuber.ChannelID).Find(&vtubers)
	if res.RowsAffected != 0 {
		c.JSON(400, gin.H{"status": "It's already created"})
		return
	}
	if res.Error != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	if err := db.Create(vtuber).Error; err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	for _, t := range requestVtuber.Types {
		vtuberType := &VtuberType{
			VtuberID: vtuber.ID,
			Types:    t,
		}
		if err := db.Create(vtuberType).Error; err != nil {
			c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
			return
		}
	}
	c.JSON(201, vtuber)
}

func RegisterVtuberJsonFile(c *gin.Context) {
	// ---------------------------
	config, err := database.NewLocalDB("user", "password", "vtuber")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// -----------------------------
	vDataBytes, err := ioutil.ReadFile("./vtuber-data/vtuber-req.json")
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	// jsonVtubers := ([]byte)(vDataBytes)
	// reqVtubers := []*RequestVtuber{}
	type reqVtubers struct {
		Data []RequestVtuber `json:"data"`
	}
	reqV := &reqVtubers{}
	if err := json.Unmarshal(vDataBytes, reqV); err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	vtuber := &Vtuber{}
	vtuberType := &VtuberType{}
	for _, vd := range reqV.Data {
		if err := db.Model(vtuber).Where("channel_id = ?", vd.ChannelID).
			Updates(&Vtuber{
				Name:        vd.Name,
				ChannelID:   vd.ChannelID,
				Affiliation: vd.Affiliation,
			}).Error; err != nil {
			c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ErrRecordDatabase))
			return
		}

		// db.Model(&VtuberType{}).Where("vtuber_id = ?",vtuber.ID).Find([]VtuberType{})
		db.Model(vtuberType).Where("vtuber_id = ?", vtuber.ID).Update()
	}
	c.JSON(200, gin.H{"status": "ok"})
}
func useIoutilReadFile(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
