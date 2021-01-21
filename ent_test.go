package main

import (
	"context"
	"fmt"
	"testing"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/enttest"
	"knowlgraph.com/ent/migrate"
	"knowlgraph.com/ent/star"
	"knowlgraph.com/ent/story"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"

	_ "github.com/mattn/go-sqlite3"
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

	softwareDevelopment := client.Node.Create().SetTag(client.Tag.Query().Where(tag.NameEQ("Software development")).FirstX(ctx)).SaveX(ctx)
	t.Log(softwareDevelopment)

	android := client.Node.Create().
		SetTag(client.Tag.Query().Where(tag.NameEQ("Android")).FirstX(ctx)).
		SetForm(softwareDevelopment).
		SetRoot(softwareDevelopment).
		SaveX(ctx)
	t.Log(android)

	iOS := client.Node.Create().
		SetTag(client.Tag.Query().Where(tag.NameEQ("iOS")).FirstX(ctx)).
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

	_response := client.Story.Create().SetState(story.StatePublic).SaveX(ctx)
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
	_content := client.Content.Create().
		SetTitle("What's KnowlGraph").
		// SetGist("Download and install the Slice Viewer. Run the Slice Viewer.").
		SetContent("A publishing platform where people can read professional and in-depth knowledge on the graph.").
		SetLang("ZH").
		SetVersion(1).
		SetStatus(content.StatusSave).
		AddTags(_tags...).
		SaveX(ctx)
	t.Log(_content)
	_story := client.Story.Create().SetState(story.StatePublic).AddVersions(_content).SaveX(ctx)
	t.Log(_story)

	_tags = client.Tag.Query().Where(tag.NameIn("ent", "Software development", "golang", "Web")).AllX(ctx)
	_content = client.Content.Create().
		SetTitle("What's Ent").
		SetGist("quick introduction of ent").
		SetContent("ent is a simple, yet powerful entity framework for Go, that makes it easy to build and maintain applications with large data-models.").
		SetLang("EN").
		SetVersion(1).
		SetStatus(content.StatusSave).
		AddTags(_tags...).
		SaveX(ctx)
	t.Log(_content)
	_story = client.Story.Create().SetState(story.StatePublic).AddVersions(_content).SaveX(ctx)
	t.Log(_story)
}

func CreateClient(t *testing.T) (context.Context, *ent.Client) {
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}

	// https://godoc.org/github.com/mattn/go-sqlite3#SQLiteDriver.Open
	// client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
	client := enttest.Open(t, "sqlite3", "file:test.db?_fk=1", opts...)

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, nil
	}
	return context.Background(), client
}
