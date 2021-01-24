package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/migrate"

	_ "github.com/mattn/go-sqlite3"
)

// Config is KnowlGraph server config info
type Config struct {
	Port  int    `flag:"server port"`
	Db    string `flag:"database path"`
	Debug bool   `flag:"Is debug mode"`
}

// DefaultLanguage is story default language
const DefaultLanguage = "en"

var ctx context.Context
var client *ent.Client
var config *Config

func init() {
	time.FixedZone("CST", 8*3600) // China Standard Timzone

	config = &Config{8080, "", false}
	goflag.Var(config)
}

func main() {
	goflag.Parse("config", "Configuration file path")

	var err error

	opts := []ent.Option{}
	if config.Debug {
		opts = append(opts, ent.Debug())
	}
	if "" == config.Db {
		client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
	} else {
		client, err = ent.Open("sqlite3", "file:"+config.Db+"?_fk=1", opts...)
	}
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	defer client.Close()

	ctx = context.Background()
	if err = client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		panic("failed to create schema: " + err.Error())
	}

	//////////// load languages.json resource //////////////
	//
	//
	////////////////////////////////////////////////////////

	jsonFile, err := os.Open("res/languages.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var values []*ent.Language
	err = json.Unmarshal(byteValue, &values)
	if err != nil {
		panic(err)
	}

	var _langCreates []*ent.LanguageCreate
	var i = 0
	for _, v := range values {
		_, err := client.Language.UpdateOneID(v.ID).SetName(v.Name).SetDirection(v.Direction).Save(ctx)
		if err != nil {
			_langCreates[i] = client.Language.Create().SetID(v.ID).SetName(v.Name).SetDirection(v.Direction)
			i++
		}
	}
	client.Language.CreateBulk(_langCreates...).SaveX(ctx)

	/////////////////////// router code ////////////////////
	//
	//
	////////////////////////////////////////////////////////

	router := gin.Default()
	v1 := router.Group("/api/v1")

	v1.PUT("/story", handle(NewStory))
	v1.PUT("/story/content", handle(PutStoryContent))
	v1.PUT("/response")
	v1.PUT("/response/content")

	v1.POST("/reaction")
	v1.POST("/star")
	v1.POST("/node")
	v1.POST("/fork")
	v1.POST("/recover")
	v1.POST("/publish", handle(PublishStory))

	v1.GET("/story", handle(GetStory))
	v1.GET("/story/responses")
	v1.GET("/story/versions")
	v1.GET("/story/reactions")
	v1.GET("/story/fork")

	v1.GET("/tags")
	v1.GET("/tag/stories")
	v1.GET("/tag/nodes")

	v1.GET("/node/stories")
	v1.GET("/node")

	v1.GET("/star/storis")
	v1.GET("/star/nodes")

	v1.GET("/recent")
	v1.GET("/search")

	router.Run(fmt.Sprint(":", config.Port))
}
