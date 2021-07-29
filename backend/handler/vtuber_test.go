package handler_test

import (
	"fmt"
	"testing"

	"github.com/1k-ct/twitter-dem/pkg/database"
	"github.com/1k-ct/vtuber-cho/backend/handler"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/xerrors"
)

func TestFitchRandVtuber(t *testing.T) {
}
func testAutomigrad() error {
	config, err := database.NewLocalDB("user", "password", "sample")
	if err != nil {
		return xerrors.New(err.Error())
	}
	db, err := config.Connect()
	if err != nil {
		return xerrors.New(err.Error())
	}
	defer db.Close()

	if err := db.DropTableIfExists(&handler.Vtuber{}, &handler.VtuberTag{}).Error; err != nil {
		return xerrors.New(err.Error())
	}
	if err := db.AutoMigrate(&handler.Vtuber{}).Error; err != nil {
		return xerrors.New(err.Error())
	}
	if err := db.AutoMigrate(&handler.VtuberTag{}).AddForeignKey("vtuber_id", "vtubers(id)", "CASCADE", "CASCADE").Error; err != nil {
		return xerrors.New(err.Error())
	}
	return nil
}
func TestMain(t *testing.T) {
	vtubers, err := handler.FetchDatabaseVtuber("hololive", "tag1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(vtubers)
}

func TestRegisterData(t *testing.T) {
	v := handler.Vtuber{
		Name:        "name",
		ChannelID:   "channel id",
		Affiliation: "honey strap",
		VtuberTags: []handler.VtuberTag{
			{Tag: "tag1"},
			// {Tag: "tag2"},
			// {Tag: "tag3"},
		},
	}
	data, err := handler.RegisterDatabaseVtuber(&v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}
