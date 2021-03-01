package main

import (
	"github.com/facebook/ent/dialect/sql"
	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func publishArticle(c *Context) error {
	var _data struct {
		Name      string   `json:"name"`
		Comment   string   `json:"comment"`
		Title     string   `json:"title"`
		Gist      string   `json:"gist"`
		Lang      string   `json:"lang"`
		Tags      []string `json:"tags"`
		ContentID int      `json:"content_id"`
		ArticleID int      `json:"article_id"`
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

	// _branche, err := client.Draft.Query().
	// 	Where(draft.And(
	// 		draft.HasUserWith(user.ID(_userID.(int))),
	// 		draft.HasSnapshotsWith(content.ID(_data.ContentID)))).
	// 	WithSnapshots(func(cq *ent.ContentQuery) {
	// 		cq.Where(content.ID(_data.ContentID))
	// 	}).
	// 	Only(ctx)

	// _content := _branche.Edges.Snapshots[0]

	_article, err := client.Article.Get(ctx, _data.ArticleID)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_content, err := client.Content.Query().
		Where(content.And(
			content.ID(_data.ContentID),
			content.HasBrancheWith(
				draft.HasUserWith(user.ID(_userID.(int)))))).
		Only(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _data.Title == "" {
		_data.Title = _content.Title
	}

	if _data.Gist == "" {
		_data.Gist = _content.Body
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
			SetComment(_data.Comment).
			SetTitle(_data.Title).
			SetGist(_data.Gist).
			SetState(_versionState).
			SetLangID(_lang.ID).
			AddTags(_tags...).
			SetContentID(_data.ContentID).
			SetArticleID(_data.ArticleID).
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
