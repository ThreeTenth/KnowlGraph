package main

import (
	"regexp"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"

	stripmd "github.com/writeas/go-strip-markdown"
)

func publishArticle(c *Context) error {
	var _data struct {
		Name    string   `json:"name"`
		Comment string   `json:"comment"`
		Title   string   `json:"title"`
		Gist    string   `json:"gist"`
		Lang    string   `json:"lang"`
		Tags    []string `json:"tags"`
		DraftID int      `json:"draft_id"`
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_lang, err := client.Language.Query().Where(language.CodeEQ(_data.Lang)).First(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_draft, err := client.Draft.Query().
		Where(draft.And(
			draft.ID(_data.DraftID),
			draft.HasUserWith(user.ID(_userID.(int))))).
		WithSnapshots(func(cq *ent.ContentQuery) {
			cq.Order(ent.Desc(draft.FieldID)).Limit(1)
		}).
		WithArticle().
		First(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _draft.State == draft.StateRead {
		return c.Unauthorized("Cannot be publish because it is read-only")
	}

	_article := _draft.Edges.Article

	if 0 == len(_draft.Edges.Snapshots) {
		return c.BadRequest("Error: Content is empty")
	}

	_content := _draft.Edges.Snapshots[0]

	if _content.Body == "" {
		return c.BadRequest("Error: Content is empty")
	}

	if _data.Gist == "" {
		_data.Gist = seo(_content.Body)
	}

	_versionState := version.StateReview
	_tagStatus := tag.StatusPublic
	if _article.Status == article.StatusPrivate {
		_versionState = version.StateRelease
		_tagStatus = tag.StatusPrivate
	}

	_tags := make([]*ent.Tag, len(_data.Tags))
	for i, v := range _data.Tags {
		_tag, err := client.Tag.Query().Where(tag.Name(v)).First(ctx)
		if err != nil {
			_tag, err = client.Tag.Create().SetName(v).SetStatus(_tagStatus).Save(ctx)
		}
		if err != nil {
			return c.InternalServerError(err.Error())
		}
		_tags[i] = _tag
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		_version, err := tx.Version.Create().
			SetName(_data.Name).
			SetComment(_data.Comment).
			SetTitle(_data.Title).
			SetGist(_data.Gist).
			SetState(_versionState).
			SetLangID(_lang.ID).
			AddTags(_tags...).
			SetContentID(_content.ID).
			SetArticleID(_article.ID).
			Save(ctx)
		if err != nil {
			return err
		}

		if _versionState == version.StateReview {
			_userIDs, err := tx.User.Query().
				Order(random()).
				Limit(10).
				Select(user.FieldID).
				Ints(ctx)
			if err != nil {
				return err
			}

			_, err = tx.RAS.Create().
				SetComment(_data.Comment).
				SetVersionID(_version.ID).
				AddVoterIDs(_userIDs...).
				Save(ctx)
			if err != nil {
				return err
			}

			_, err = _draft.Update().SetState(draft.StateRead).Save(ctx)

			if err != nil {
				return err
			}
		} else {
			_, err = tx.Asset.Create().SetArticle(_article).SetUserID(_userID.(int)).SetStatus(asset.StatusSelf).Save(ctx)

			if err != nil {
				return err
			}

			err = tx.Draft.DeleteOne(_draft).Exec(ctx)

			if err != nil {
				return err
			}
		}

		return c.Ok(_version)
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return nil
}

func random() ent.OrderFunc {
	return func(s *sql.Selector, check func(string) bool) {
		f := "random()"
		if check(f) {
			s.OrderBy(f)
		} else {
			s.AddError(errors.Errorf("invalid field %q for ordering", f))
		}
	}
}

func seo(md string) string {
	// Golang string is saved in UTF-8 format.
	// For languages that use multiple bytes to express one character,
	// such as Chinese characters, when substring the string,
	// invalid character encoding will appear. We need to convert string to []rune type.
	// Rune is stored in unicode encoding, it supports characters with long byte encoding.
	_rune := []rune(stripmd.Strip(md))

	var _sub string
	if len(_rune) <= 200 {
		_sub = _sub + " " + string(_rune)
	} else {
		_sub = _sub + " " + string(_rune[:120]) + "..."
	}

	re := regexp.MustCompile(`\r?\n`)
	_sub = re.ReplaceAllString(_sub, " ")
	return strings.Trim(_sub, " ")
}
