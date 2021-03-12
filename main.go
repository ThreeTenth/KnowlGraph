package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/packr/v2"
	"golang.org/x/sync/errgroup"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/migrate"

	"github.com/facebook/ent/dialect"

	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Config is KnowlGraph server config info
type Config struct {
	Port  int    `flag:"server port"`
	Pg    string `flag:"postgresql database source name, Please enter in the format: user:password@127.0.0.1:5432/database"`
	Redis string `flag:"redis source name, Please enter in the format: redis://localhost:6379/<db>"`
	Log   string `flag:"logcat file path"`
	Debug bool   `flag:"Is debug mode"`
	Gci   string `flag:"GitHub client ID, see https://docs.github.com/en/developers/apps/creating-an-oauth-app"`
	Gcs   string `flag:"GitHub client secrets"`
	Ssd   string `flag:"static server domain"`
	Rad   string `flag:"restful api domain"`
}

var ctx context.Context
var client *ent.Client
var rdb *redis.Client
var config *Config
var langs []*ent.Language
var defaultLang *ent.Language

func init() {
	time.FixedZone("CST", 8*3600) // China Standard Timzone

	config = &Config{Port: 8080, Ssd: "http://localhost:20012", Rad: "http://localhost:20011"}
	goflag.Var(config)
}

func openPostgreSQL() {
	db, err := sql.Open("pgx", "postgresql://"+config.Pg)
	if err != nil {
		panic("open postgresql failed: " + err.Error())
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	opts := []ent.Option{
		ent.Driver(drv),
	}
	if config.Debug {
		opts = append(opts, ent.Debug())
	}
	client = ent.NewClient(opts...)

	ctx = context.Background()
	if err = client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		panic("failed to create schema: " + err.Error())
	}
	// sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt buster-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
}

func openRedis() {
	opt, err := redis.ParseURL(config.Redis)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)
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

	err = json.Unmarshal(bs, &langs)
	if err != nil {
		panic(err)
	}

	var _langCreates []*ent.LanguageCreate
	var i = 0
	for _, v := range langs {
		_langID, _ := client.Language.Update().Where(language.CodeEQ(v.Code)).SetName(v.Name).SetDirection(v.Direction).Save(ctx)
		if _langID == 0 {
			_langCreates = append(_langCreates, client.Language.Create().SetCode(v.Code).SetName(v.Name).SetDirection(v.Direction))
			i++
		}
	}
	client.Language.CreateBulk(_langCreates...).SaveX(ctx)

	langs = client.Language.Query().Order(ent.Asc(language.FieldID)).AllX(ctx)
	defaultLang = getDefaultLanguageIfNotFound(DefaultLanguage)
}

func getLanguage(code string) *ent.Language {
	for _, v := range langs {
		if v.Code == code {
			return v
		}
	}

	return nil
}

func getDefaultLanguageIfNotFound(code string) *ent.Language {
	for _, v := range langs {
		if v.Code == code {
			return v
		}
	}

	for _, v := range langs {
		if v.Code == DefaultLanguage {
			return v
		}
	}
	panic("server err: no default language.")
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

// var articleExist validator.Func = func(fl validator.FieldLevel) bool {
// 	date, ok := fl.Field().Interface().(int)
// 	if ok {
// 		ok, err := client.Article.Query().Where(article.IDEQ(date)).Exist(ctx)
// 		if !ok || err != nil {
// 			return false
// 		}
// 	}
// 	return true
// }

func main() {
	goflag.Parse("config", "Configuration file path")

	openPostgreSQL()
	openRedis()
	loadLauguages()

	/////////////////////// router code ////////////////////
	//
	//
	////////////////////////////////////////////////////////

	// write the logs to file and console at the same time
	if f, err := os.Create(config.Log); err == nil {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	server01 := &http.Server{
		Addr:         ":20010",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":20011",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server03 := &http.Server{
		Addr:         ":20012",
		Handler:      router03(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	var g errgroup.Group

	g.Go(func() error {
		err := server01.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		return err
	})

	g.Go(func() error {
		err := server02.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		return err
	})

	g.Go(func() error {
		err := server03.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func cors(c *gin.Context) {
	// todo
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func router01() http.Handler {
	router := gin.Default()
	loadTemplates(router)

	router.Use(cors)

	router.GET("/favicon.ico", getFavicon)

	router.GET("/", authentication, html(index))
	router.GET("/drafts/*id", authentication, html(index))
	router.GET("/about", authentication, html(index))
	router.GET("/my", authentication, html(index))
	router.GET("/new/article", authentication, html(index))
	router.GET("/p/:id/edit", authentication, html(index))

	router.GET("/signout", deauthorize, handle(signout))

	join := router.Group("/user/join")
	join.GET("/github", handle(joinGithub))

	return router
}

func router02() http.Handler {
	router := gin.Default()
	router.Use(cors)
	router.GET("/favicon.ico", getFavicon)

	v1 := router.Group("/v1")

	v1.PUT("/article", authorizeRequired, handle(putArticleNew))
	v1.PUT("/article/content", authorizeRequired, handle(putArticleContent))

	v1.PUT("/publish/article", authorizeRequired, handle(publishArticle))

	v1.GET("/drafts", authorizeRequired, handle(getArticleDrafts))
	v1.GET("/user/articles", authorizeRequired, handle(getUserArticles))

	return router
}

func router03() http.Handler {
	router := gin.Default()
	router.Use(cors)

	router.GET("/favicon.ico", getFavicon)
	router.GET("/static/*paths", getStaticServerFiles)
	router.GET("/theme/theme.js", getStaticTheme)
	router.GET("/theme/theme@:id.js", getStaticTheme)
	router.GET("/lang/:id", getStaticLang)

	return router
}
