package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
)

const FILE_PATH = "./../vtuber-data/vtuber-req.json"

type Item struct {
	Types []string `json:"types"`
}
type Data struct {
	VtuberData []Vtuber `json:"data"`
}
type Vtuber struct {
	Name        string   `json:"name"`
	ChannelID   string   `json:"channel_id"`
	Affiliation string   `json:"affiliation"`
	Types       []string `json:"types"`
}

type UserResponse struct {
	RequestMethod string      `json:"RequestMethod"`
	Result        interface{} `json:"Result"`
}

func HandlerRandItem(c *gin.Context) {
	result := &UserResponse{}
	affiliationParam := c.Param("affiliations")
	typeParam := c.Param("types")

	if affiliationParam == "" && typeParam == "" {
		result.Result = errors.New("the url is incorrect")
		c.JSON(400, result)
		return
	}

	res, err := FindItems(FILE_PATH, affiliationParam, typeParam)
	if err != nil {
		c.JSON(500, res)
		return
	}
	var response struct {
		Name string
		URL  string
	}
	response.Name = res.Name
	response.URL = fmt.Sprintf("https://www.youtube.com/channel/%v", res.ChannelID)

	result.RequestMethod = "GET"
	result.Result = response

	c.JSON(200, result)
}
func HandlerItemSearch(c *gin.Context) {
	result := &UserResponse{}
	channelID := c.Param("channel")
	v, err := FitchItem(FILE_PATH, channelID)
	if err != nil {
		result.Result = errors.New("server error: json file error")
		c.JSON(500, result)
		return
	}
	result.Result = v
	c.JSON(200, result)
}
func FindItems(file, targetAffiliation, targetType string) (*Vtuber, error) {
	item := []int{}
	dataCou := gojsonq.New().File(file).From("data").Count()

	for i := 0; i < dataCou; i++ {
		typesForm := fmt.Sprintf("data.[%v].types", i)
		res, err := gojsonq.New().File(file).From(typesForm).GetR()
		if err != nil {
			return nil, err
		}
		r, _ := res.StringSlice()

		affForm := fmt.Sprintf("data.[%v].affiliation", i)
		affiliations, err := gojsonq.New().File(file).From(affForm).GetR()
		if err != nil {
			return nil, err
		}
		affiliation, _ := affiliations.String()

		if targetAffiliation == affiliation && stringInSlice(targetType, r) {
			item = append(item, i)
		}
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(item))

	resForm := fmt.Sprintf("data.[%v]", item[n])
	res := gojsonq.New().File(file).From(resForm).Get()
	// ress := res.(map[string]interface{})
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	v := &Vtuber{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}
func marshalVtuber(obj interface{}) (*Vtuber, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	v := &Vtuber{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}
func FitchItem(file, channelID string) (interface{}, error) {
	res := gojsonq.New().File(file).From("data").
		Where("channel_id", "=", channelID).Get()
	return res, nil
}
func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}

	return false
}
