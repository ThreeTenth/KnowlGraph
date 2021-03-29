package main

import (
	"regexp"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
	"knowlgraph.com/ent/word"

	stripmd "github.com/writeas/go-strip-markdown"
)

func publishArticle(c *Context) error {
	var _data struct {
		Name     string   `json:"name"`
		Comment  string   `json:"comment"`
		Title    string   `json:"title"`
		Gist     string   `json:"gist"`
		Lang     string   `json:"lang"`
		Keywords []string `json:"keywords"`
		DraftID  int      `json:"draft_id"`
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	// todo 校验 lang code 的准确性

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
	if _article.Status == article.StatusPrivate {
		_versionState = version.StateRelease
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {

		// todo publish successed, the words has set public status
		_words := make([]*ent.Word, len(_data.Keywords))
		for i, v := range _data.Keywords {
			_word, err := tx.Word.Query().Where(word.Name(v)).First(ctx)
			if err != nil {
				_word, err = tx.Word.Create().SetName(v).SetStatus(word.StatusPrivate).Save(ctx)
			}
			if err != nil {
				return err
			}
			_words[i] = _word
		}

		_version, err := tx.Version.Create().
			SetName(_data.Name).
			SetComment(_data.Comment).
			SetTitle(_data.Title).
			SetGist(_data.Gist).
			SetState(_versionState).
			SetLang(_data.Lang).
			AddKeywords(_words...).
			SetContentID(_content.ID).
			SetArticleID(_article.ID).
			Save(ctx)
		if err != nil {
			return err
		}

		_version.Edges.Content = _content
		_version.Edges.Article = _article
		_version.Edges.Keywords = _words

		if _versionState == version.StateRelease {
			if err = updateUserAsset(tx, _userID.(int), asset.StatusSelf, _version); err != nil {
				return err
			}

			return c.Ok(_version)
		}

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

		return c.NoContent()
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

func updateUserAsset(tx *ent.Tx, userID int, status asset.Status, vers *ent.Version) error {
	_, err := tx.Asset.Query().
		Where(asset.HasArticleWith(
			article.ID(vers.Edges.Article.ID))).
		Only(ctx)

	if err != nil {
		_, err = tx.Asset.Create().
			SetArticle(vers.Edges.Article).
			SetUserID(userID).
			SetStatus(status).
			Save(ctx)

		if err != nil {
			return err
		}
	}

	_, err = tx.User.Update().
		Where(user.ID(userID)).
		AddWords(vers.Edges.Keywords...).
		Save(ctx)

	if err != nil {
		return err
	}

	_, err = tx.Language.Create().
		SetCode(vers.Lang).
		SetUserID(userID).
		Save(ctx)

	if err != nil {
		return err
	}

	count, err := tx.Draft.Delete().
		Where(draft.HasSnapshotsWith(
			content.ID(vers.Edges.Content.ID))).
		Exec(ctx)
	if 0 == count {
		return errors.New("No draft")
	}

	return err
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
