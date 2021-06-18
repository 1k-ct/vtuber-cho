package handler_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestFitchRandVtuber(t *testing.T) {
	// var c *gin.Context
	// handler.FitchRandVtuber(c)
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(os.Getenv("API_KEY"))
}
