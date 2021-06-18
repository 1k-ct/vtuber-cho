package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/1k-ct/twitter-dem/pkg/appErrors"
	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/1k-ct/twitter-dem/pkg/jwtToken"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type LoginInfo struct {
	APIKey string
}
type ResponseJson struct {
	Status  string
	Message string
	Data    interface{}
}

func FitchJwt(c *gin.Context) {
	// LoginInfo := &LoginInfo{}
	// if err := c.BindJSON(&LoginInfo); err != nil {
	// 	c.JSON(400, appErrors.ErrMeatdataMsg(err, appErrors.ErrorJSON))
	// 	return
	// }

	// api key の取得
	// apikey := c.Param("api_key")

	// if apikey == "" {
	// 	c.JSON(400, &ResponseJson{
	// 		Status:  "400",
	// 		Message: "The api key is required",
	// 		Data:    nil,
	// 	})
	// }
	// // 取得したapi key から、hash　を生成する
	// hash, err := bcrypt.GenerateFromPassword([]byte(apikey), 10)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest,
	// 		&ResponseJson{
	// 			Status:  "400",
	// 			Message: "There is a problem with the api_key.",
	// 			Data:    nil,
	// 		})
	// }
	// apikey = string(hash)

	// .envにあるapi key と　hash を比較

	// loginInfo := &LoginInfo{}
	// if err := c.ShouldBindJSON(&loginInfo); err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, appErrors.ErrMeatdataMsg(err, appErrors.ErrorJSON))
	// 	return
	// }

	apiKey := c.Request.Header.Get("api-key")
	// user, err := ah.accountUseCase.FindByID(u.ID)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		c.JSON(http.StatusUnauthorized, appErrors.ErrMeatdataMsg(err, appErrors.ErrRecordDatabase))
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, appErrors.ServerError)
	// 	return
	// }
	if err := godotenv.Load(); err != nil {
		c.JSON(500, appErrors.ErrMeatdataMsg(err, appErrors.ServerError))
		return
	}
	// log.Println(apiKey)
	log.Println("api-key is", os.Getenv("API_KEY"))
	if os.Getenv("API_KEY") != apiKey {
		c.JSON(400, gin.H{"status": "Invalid key"})
		return
	}
	ts, err := jwtToken.CreateToken("admin", "admin", os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_KEY"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, appErrors.ErrNotCreateToken)
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)

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

	type RequestVtuber struct {
		Name        string   `json:"name"`
		ChannelID   string   `json:"channel_id"`
		Affiliation string   `json:"affiliation"`
		Types       []string `json:"types"`
	}
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
		c.JSON(400, gin.H{"status": "すでにあります。"})
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
	c.JSON(200, vtuber)
	// TODO:
	// if user.UserName == "" {
	// 	c.JSON(400, gin.H{"msg": "UserNameは必須です。"})
	// 	return
	// }
	// if user.Password == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"msg": "Passwordは必須です。"})
	// }
}
