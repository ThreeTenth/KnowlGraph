package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/packr/v2"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/story"

	"github.com/facebook/ent/dialect"

	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Config is KnowlGraph server config info
type Config struct {
	Port  int    `flag:"server port"`
	Pg    string `flag:"postgresql database source name, Please enter in the format: user:password@127.0.0.1:5432/database"`
	Log   string `flag:"logcat file path"`
	Debug bool   `flag:"Is debug mode"`
	Gci   string `flag:"GitHub client ID, see https://docs.github.com/en/developers/apps/creating-an-oauth-app"`
	Gcs   string `flag:"GitHub client secrets"`
}

var ctx context.Context
var client *ent.Client
var config *Config

func init() {
	time.FixedZone("CST", 8*3600) // China Standard Timzone

	config = &Config{Port: 8080}
	goflag.Var(config)
}

func openPostgreSQL() {
	db, err := sql.Open("pgx", "postgresql://"+config.Pg)
	if err != nil {
		panic("open postgresql failed: " + err.Error())
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	client = ent.NewClient(ent.Driver(drv))

	ctx = context.Background()
	if err = client.Schema.Create(ctx); err != nil {
		panic("failed to create schema: " + err.Error())
	}
}

func loadLauguages() {
	//////////// load languages.json resource //////////////
	//
	//
	////////////////////////////////////////////////////////

	box := packr.NewBox("./res")
	bs, err := box.Find("languages.json")
	if err != nil {
		panic(err)
	}

	var values []*ent.Language
	err = json.Unmarshal(bs, &values)
	if err != nil {
		panic(err)
	}

	var _langCreates []*ent.LanguageCreate
	var i = 0
	for _, v := range values {
		_, err := client.Language.UpdateOneID(v.ID).SetName(v.Name).SetDirection(v.Direction).Save(ctx)
		if err != nil {
			_langCreates = append(_langCreates, client.Language.Create().SetID(v.ID).SetName(v.Name).SetDirection(v.Direction))
			i++
		}
	}
	client.Language.CreateBulk(_langCreates...).SaveX(ctx)
}

func loadTemplates(router *gin.Engine) {
	tmpl := template.New("user")
	box := packr.NewBox("./tpl")

	for _, v := range box.List() {
		fgm := tmpl.New(v)
		data, _ := box.FindString(v)
		fgm.Parse(data)
	}

	router.SetHTMLTemplate(tmpl)
}

var storyExist validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(int)
	if ok {
		ok, err := client.Story.Query().Where(story.IDEQ(date)).Exist(ctx)
		if !ok || err != nil {
			return false
		}
	}
	return true
}

func main() {
	goflag.Parse("config", "Configuration file path")

	openPostgreSQL()
	loadLauguages()

	/////////////////////// router code ////////////////////
	//
	//
	////////////////////////////////////////////////////////

	router := gin.Default()

	// write the logs to file and console at the same time
	if f, err := os.Create(config.Log); err == nil {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	// Custom Validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("storyExist", storyExist)
	}

	web := router.Group("/")
	{
		loadTemplates(router)

		web.GET("/", authentication, handle(index))
		web.GET("/signout", deauthorize, handle(signout))

		join := web.Group("/user/join")
		join.GET("/github", handle(joinGithub))
	}

	web.StaticFS("/favicon", packr.NewBox("./res/favicon"))

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
