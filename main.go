package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/packr"
	"golang.org/x/sync/errgroup"
	"knowlgraph.com/ent"
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
	Build bool   `flag:"Is build mode. This mode will start pre-compiling static resource files to build directory"`
	Gci   string `flag:"GitHub client ID, see https://docs.github.com/en/developers/apps/creating-an-oauth-app"`
	Gcs   string `flag:"GitHub client secrets"`
	Ssd   string `flag:"static server domain"`
	Rad   string `flag:"restful api domain"`
}

var ctx context.Context
var client *ent.Client
var rdb *redis.Client
var config *Config
var db *sql.DB

func init() {
	fmt.Printf("The current version: v%v, the version name: %v", Version, VersionName)
	time.FixedZone("CST", 8*3600) // China Standard Timzone

	config = &Config{Port: 8080, Ssd: "http://localhost:20012", Rad: "http://localhost:20011"}
	goflag.Var(config)
}

func panicIfErrNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func openPostgreSQL() {
	var err error
	db, err = sql.Open("pgx", "postgresql://"+config.Pg)
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
	err = client.Schema.Create(ctx,
		migrate.WithGlobalUniqueID(true),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true))
	if err != nil {
		panic("failed to create schema: " + err.Error())
	}
}

func openRedis() {
	opt, err := redis.ParseURL(config.Redis)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)

	// 获取当前缓存的静态资源版本名称，
	// 并与应用程序（Knowlgraph 服务端程序）的版本进行比较，
	// 如果一致，则返回，
	// 如果不一致，则删除所有静态资源的缓存。
	// 缓存静态资源 @see getJsdeliverFile(_filepath)
	_versionName, _ := rdb.Get(ctx, SFVersion).Result()
	if _versionName == VersionName {
		return
	}

	err = rdb.Set(ctx, SFVersion, VersionName, ExpireTimeToken).Err()
	panicIfErrNotNil(err)
	sfkeys, err := rdb.Keys(ctx, SF+"*").Result()
	panicIfErrNotNil(err)
	if 0 == len(sfkeys) {
		return
	}

	_, err = rdb.Del(ctx, sfkeys...).Result()
	panicIfErrNotNil(err)
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
	config.Ssd += "/" + VersionName

	// 开发阶段时使用本地资源,
	// 正式版本使用编译后的资源
	if config.Build {
		buildStaticFile()
		return
	}

	openPostgreSQL()
	openRedis()

	defer db.Close()

	startServer()
}

func startServer() {
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
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, HEAD")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func router01() http.Handler {
	router := gin.Default()
	loadTemplates(router)

	if config.Debug {
		router.Use(cors)
	}

	router.GET("/favicon.ico", getFavicon)

	router.GET("/", authentication, html(index))
	router.GET("/drafts", authentication, html(index))
	router.GET("/archive", authentication, html(index))
	router.GET("/archive/*path", authentication, html(index))
	router.GET("/about", authentication, html(index))
	router.GET("/my", authentication, html(index))
	router.GET("/new/article", authentication, html(index))
	router.GET("/d/*path", authentication, html(index))
	router.GET("/p/:id", authentication, html(articleHTML))
	router.GET("/p/:id/*path", authentication, html(articleHTML))
	router.GET("/settings/*path", authentication, html(index))

	router.GET("/login", authentication, html(index))
	router.GET("/signout", deauthorize, handle(signout))

	join := router.Group("/user/join")
	join.GET("/github", handle(joinGithub))

	return router
}

func router02() http.Handler {
	router := gin.Default()

	if config.Debug {
		router.Use(cors)
	}

	router.GET("/favicon.ico", getFavicon)

	v1 := router.Group("/v1")

	v1.PUT("/article", authorizeRequired, handle(putArticleNew))
	v1.PUT("/article/content", authorizeRequired, handle(putArticleContent))
	v1.PUT("/article/edit", authorizeRequired, handle(editArticle))

	v1.PUT("/publish/article", authorizeRequired, handle(publishArticle))

	v1.GET("/articles", authentication, handle(getArticles))
	v1.GET("/words", authentication, handle(getKeywords))
	v1.GET("/drafts", authorizeRequired, handle(getDrafts))
	v1.GET("/draft", authorizeRequired, handle(getDraft))
	v1.GET("/article", authentication, handle(getArticle))

	v1.GET("/node", handle(getNode))
	v1.PUT("/node", authorizeRequired, handle(putNode))
	v1.GET("/node/articles", handle(getNodeArticles))

	v1.GET("/archives", authorizeRequired, handle(getArchives))
	v1.PUT("/archive", authorizeRequired, handle(putArchive))

	v1.GET("/assets", authorizeRequired, handle(getAssets))
	v1.PUT("/asset", authorizeRequired, handle(putAsset))
	v1.DELETE("/asset", authorizeRequired, handle(deleteAsset))

	v1.PUT("reaction", authentication, handle(putReaction))

	v1.GET("/vote", authorizeRequired, handle(getVote))
	v1.POST("/vote", authorizeRequired, handle(postVote))

	v1.GET("/word", handle(getWord))
	v1.GET("/word/articles", handle(getKeywordArticles))
	v1.GET("/word/nodes", handle(getWordNodes))

	v1.GET("/account/challenge")
	v1.PUT("/account/create")
	v1.POST("/account/request")
	v1.POST("/account/anthn")
	v1.POST("/account/sync")

	return router
}

func router03() http.Handler {
	router := gin.Default()

	if config.Debug {
		router.Use(cors)
	}

	group := router.Group("/" + VersionName)

	if config.Debug {
		group.Static("/code-of-conduct", "./res/code-of-conduct")
		group.Static("/strings", "./res/strings")
		group.GET("/__app.js", func(c *gin.Context) {

			languagesFormat := "// %v, %v \n const languages = %v"
			defaultStringsFormat := "// %v, %v \n const defaultLang = %v"

			c.Status(http.StatusOK)
			c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")

			CopyFile(c.Writer, "./res/languages.json", languagesFormat)
			CopyFile(c.Writer, "./res/strings/strings-en.json", defaultStringsFormat)
			CopyDir(c.Writer, "./res/js/components", "")
			CopyDir(c.Writer, "./res/js/routers", "")
			CopyDir(c.Writer, "./res/js/utils", "")
			CopyFile(c.Writer, "./res/js/app.js", "")
		})
		group.GET("/__default_theme.js", func(c *gin.Context) {

			themeFormat := "// %v \n const %v = `%v`"

			c.Status(http.StatusOK)
			c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")

			CopyDir(c.Writer, "./res/theme", themeFormat)
		})
		group.GET("/__main.css", func(c *gin.Context) {

			c.Status(http.StatusOK)
			c.Writer.Header().Set("Content-Type", "text/css; charset=utf-8")

			CopyDir(c.Writer, "./res/css", "")
		})
	} else {
		group.GET("/*filepath", getStaticFile)
		group.HEAD("/*filepath", getStaticFile)
	}

	router.GET("/favicon.ico", getFavicon)
	// router.GET("/static/*paths", getStaticServerFiles)
	// router.GET("/theme/theme.js", getStaticTheme)
	// router.GET("/theme/theme@:id.js", getStaticTheme)
	// router.GET("/lang/:id", getStaticLang)

	return router
}
