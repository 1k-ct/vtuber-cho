package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/1k-ct/vtuber-cho/handler"
	"github.com/joho/godotenv"
)

func TestGenerateFromPassword(t *testing.T) {

}
func TestJsonReadFile(t *testing.T) {
	vDataBytes, err := ioutil.ReadFile("../vtuber-data/vtuber-req.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(vDataBytes)
}

type requestVtubers struct {
	Data []handler.RequestVtuber `json:"data"`
}

func TestReqestVtuberJsonFile(t *testing.T) {
	reqVtubers := &requestVtubers{}
	vDataBytes, err := ioutil.ReadFile("../vtuber-data/vtuber-req.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(vDataBytes, reqVtubers); err != nil {
		t.Fatal(err)
	}
	url := "http://localhost:8080/api/v1/vtubers"
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	accessToken := os.Getenv("ACCESS_TOKEN")
	log.Println(accessToken)
	if accessToken == "" {
		t.Fatal("no found access_token")
	}
	for _, v := range reqVtubers.Data {
		vData, _ := json.Marshal(v)
		body, err := testPostVtuber(url, accessToken, vData)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(body)
	}
}
func testPostVtuber(url, accessToken string, jsonStr []byte) (string, error) {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(body), err
}
func TestCreateDatabase(t *testing.T) {
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
	if err := db.AutoMigrate(handler.Vtuber{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(handler.VtuberType{}).
		AddForeignKey("vtuber_id", "vtubers(id)", "CASCADE", "CASCADE").Error; err != nil {
		t.Fatal(err)

	}
}
