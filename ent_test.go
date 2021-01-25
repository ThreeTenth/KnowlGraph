package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/star"
	"knowlgraph.com/ent/story"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	// _ "github.com/mattn/go-sqlite3"
)

func TestTag(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	names := []string{"iOS", "Android", "Software development", "Web", "iOS development", "ent", "golang"}
	tagCreaters := make([]*ent.TagCreate, len(names))
	for i, name := range names {
		tagCreaters[i] = client.Tag.Create().SetName(name)
	}

	tags := client.Tag.CreateBulk(tagCreaters...).SaveX(ctx)

	t.Log(tags)
}

func TestQueryTag(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_tag, err := client.Tag.Query().Where(tag.NameEqualFold("ios")).First(ctx)
	t.Log(_tag, err)
}

func TestUser(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	names := []string{"Tom", "Tiwen", "GaoZhao"}
	userCreaters := make([]*ent.UserCreate, len(names))
	for i, name := range names {
		userCreaters[i] = client.User.Create().SetName(name)
	}

	users := client.User.CreateBulk(userCreaters...).SaveX(ctx)

	t.Log(users)
}

func TestNode(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	softwareDevelopment := client.Node.Create().SetLevel(0).SetTag(client.Tag.Query().Where(tag.NameEQ("Software development")).FirstX(ctx)).SaveX(ctx)
	t.Log(softwareDevelopment)

	android := client.Node.Create().
		SetTag(client.Tag.Query().Where(tag.NameEQ("Android")).FirstX(ctx)).
		SetLevel(softwareDevelopment.Level + 1).
		SetForm(softwareDevelopment).
		SetRoot(softwareDevelopment).
		SaveX(ctx)
	t.Log(android)

	iOS := client.Node.Create().
		SetTag(client.Tag.Query().Where(tag.NameEQ("iOS")).FirstX(ctx)).
		SetLevel(softwareDevelopment.Level + 1).
		SetForm(softwareDevelopment).
		SetRoot(softwareDevelopment).
		SaveX(ctx)
	t.Log(iOS)
}

func TestStar(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	tom := client.User.Query().Where(user.NameEQ("Tom")).OnlyX(ctx)
	_story := client.Content.Query().Where(content.TitleContains("KnowlGraph")).QueryStory().FirstX(ctx)
	_starStory := client.Star.Create().
		SetStatus(star.StatusStar).
		SetOwner(tom).
		SetStory(_story).
		SaveX(ctx)

	_node := client.Tag.Query().Where(tag.NameEQ("Android")).QueryNodes().FirstX(ctx)
	_starNode := client.Star.Create().
		SetStatus(star.StatusStar).
		SetOwner(tom).
		SetNode(_node).
		SaveX(ctx)

	t.Log(tom)
	t.Log(_story)
	t.Log(_starStory)
	t.Log(_starNode)
}

func TestQueryStar(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_stars := client.User.Query().Where(user.NameEQ("Tom"))
	_starNodes := _stars.QueryStars().Where(star.HasNode()).AllX(ctx)
	_starStories := _stars.QueryStars().Where(star.HasStory()).AllX(ctx)

	t.Log(_stars)
	t.Log(_starNodes)
	t.Log(_starStories)
}

func TestArchive(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_stars := client.User.Query().Where(user.NameEQ("Tom"))
	_starNode := _stars.QueryStars().Where(star.HasNode()).FirstX(ctx)
	_starStory := _stars.QueryStars().Where(star.HasStory()).FirstX(ctx)
	t.Log(_starNode)
	t.Log(_starStory)

	_starNode = _starNode.Update().AddStories(_starStory).SaveX(ctx)
	t.Log(_starNode)
}

func TestQuote(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_content := client.Content.Query().Where(content.TitleContains("KnowlGraph")).FirstX(ctx)
	_story := _content.QueryStory().FirstX(ctx)

	_response := client.Story.Create().SetStatus(story.StatusPublic).SaveX(ctx)
	t.Log(_response)

	_quote := client.Quote.Create().
		SetText("read professional and in-depth knowledge").
		SetContext("A publishing platform where people can read professional and in-depth knowledge on the graph.").
		SetSource(fmt.Sprintf("db://story/%d/%d", _story.ID, _content.ID)).
		SetResponse(_response).
		SaveX(ctx)
	t.Log(_quote)
}

func TestReaction(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_story := client.Content.Query().Where(content.TitleContains("KnowlGraph")).QueryStory().FirstX(ctx)
	t.Log(_story)

	_reaction := client.Reaction.Create().SetName("+1").SetEmoji("😀").SetCount(10).AddStorise(_story).SaveX(ctx)
	t.Log(_reaction)
}

func TestStory(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	_tags := client.Tag.Query().Where(tag.NameIn("Web", "Software development")).AllX(ctx)
	_story := client.Story.Create().SetStatus(story.StatusPublic).SaveX(ctx)
	t.Log(_story)
	_content := client.Content.Create().
		SetTitle("什么是知识图谱").
		// SetGist("Download and install the Slice Viewer. Run the Slice Viewer.").
		SetContent("一个发布文章即可获取可持续收入的，可以阅读专业有深度的知识，建立丰富的知识图谱的一个地方").
		SetVersion(5).
		SetLangID("zh").
		SetStoryID(_story.ID).
		AddTags(_tags...).
		SaveX(ctx)
	t.Log(_content)

	// en := client.Language.Query().Where(language.IDEQ("en")).FirstIDX(ctx)
	// _tags = client.Tag.Query().Where(tag.NameIn("ent", "Software development", "golang", "Web")).AllX(ctx)
	// _story = client.Story.Create().SetState(story.StatePublic).SaveX(ctx)
	// t.Log(_story)
	// _content = client.Content.Create().
	// 	SetTitle("What's Ent").
	// 	SetGist("quick introduction of ent").
	// 	SetContent("ent is a simple, yet powerful entity framework for Go, that makes it easy to build and maintain applications with large data-models.").
	// 	SetLangID(en).
	// 	SetStory(_story).
	// 	SetVersion(1).
	// 	SetVersionName("v1.6.3").
	// 	AddTags(_tags...).
	// 	SaveX(ctx)
	// t.Log(_content)
}

func TestQueryStory(t *testing.T) {
	ctx, client := CreateClient(t)
	defer client.Close()

	var columns []struct {
		Lang string `json:"content_lang"`
		Max  int    `json:"max"`
	}

	client.Story.GetX(ctx, 25769803777).
		QueryVersions().
		// Where(content.ID(25769803777)).
		GroupBy(content.LangColumn).
		Aggregate(ent.Max(content.FieldVersion)).
		ScanX(ctx, &columns)
	t.Log(columns)
}

func CreateClient(t *testing.T) (context.Context, *ent.Client) {
	db, err := sql.Open("pgx", "postgresql://postgres:123456@127.0.0.1:5432/knowlgraph")
	if err != nil {
		panic("open postgresql failed: " + err.Error())
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	client := ent.NewClient(ent.Driver(drv))

	ctx := context.Background()
	if err = client.Schema.Create(ctx); err != nil {
		panic("failed to create schema: " + err.Error())
	}

	jsonFile, err := os.Open("./res/languages.json")
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

	for _, v := range values {
		_, err := client.Language.UpdateOneID(v.ID).SetName(v.Name).SetDirection(v.Direction).Save(ctx)
		if err != nil {
			client.Language.Create().SetID(v.ID).SetName(v.Name).SetDirection(v.Direction).SaveX(ctx)
		}
	}

	return ctx, client
}
