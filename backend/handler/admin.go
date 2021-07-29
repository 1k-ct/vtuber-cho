package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/1k-ct/twitter-dem/pkg/appErrors"
	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/1k-ct/twitter-dem/pkg/jwtToken"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DatabaseConnection prod server
func DatabaseConnection() (*gorm.DB, error) {
	config := database.ConfigList{
		DbDriverName:   "mysql",
		DbName:         os.Getenv("DB_NAME"),
		DbUserName:     os.Getenv("DB_USER"),
		DbUserPassword: os.Getenv("DB_PASS"),
		DbHost:         os.Getenv("DB_ADDRESS"),
		DbPort:         os.Getenv("DB_PORT"),
	}
	CONNECT := config.DbUserName + ":" + config.DbUserPassword + "@tcp(" + config.DbHost + ":3306)/" + config.DbName + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(config.DbDriverName, CONNECT)
	if err != nil {
		return db, err
	}
	return db, nil
}

// DatabaseConnection localhost server
// func DatabaseConnection() (*gorm.DB, error) {
// 	// ---------------------------
// 	config, err := database.NewLocalDB("user", "password", "vtuber")
// 	if err != nil {
// 		return nil, err
// 	}
// 	db, err := config.Connect()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// 	// -----------------------------
// }
func Getenv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return value, errors.New("not found env key")
	}
	return value, nil
}
func FitchJwt(c *gin.Context) {
	apiKey := c.Query("API_KEY")

	if os.Getenv("API_KEY") != apiKey {
		c.JSON(401, gin.H{"status": "Invalid key"})
		return
	}
	userID := os.Getenv("USER_ID")
	userName := os.Getenv("USER_NAME")
	secretKey := os.Getenv("SECRET_KEY")
	refreshKey := os.Getenv("REFRESH_KEY")
	if userID == "" || userName == "" || secretKey == "" || refreshKey == "" {
		err := errors.New("not found env key")
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}

	ts, err := jwtToken.CreateToken(userID, userName, secretKey, refreshKey)
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
	Name        string
	ChannelID   string
	Affiliation string
	VtuberTags  []string
}

func RegisterVtuber(c *gin.Context) {
	db, err := DatabaseConnection()
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	defer db.Close()

	reqVtuber := &RequestVtuber{}

	if err := c.BindJSON(&reqVtuber); err != nil {
		c.JSON(400, appErrors.ErrMeatdataMsg(err, appErrors.ErrorJSON))
		return
	}
	tag := []VtuberTag{}
	for _, v := range reqVtuber.VtuberTags {
		t := VtuberTag{Tag: v}
		tag = append(tag, t)
	}
	vtuber := &Vtuber{
		Name:        reqVtuber.Name,
		ChannelID:   reqVtuber.ChannelID,
		Affiliation: reqVtuber.Affiliation,
		VtuberTags:  tag,
	}

	v, err := RegisterDatabaseVtuber(vtuber)
	if err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}

	c.JSON(201, v)
}
