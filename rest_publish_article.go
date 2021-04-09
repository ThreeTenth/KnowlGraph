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
	"knowlgraph.com/ent/userword"
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

		// Randomly select up to 10 users as voters
		_voterIDs, err := tx.User.Query().
			Order(random()).
			Limit(10).
			Select(user.FieldID).
			Ints(ctx)
		if err != nil {
			return err
		}

		// Open a random anonymous space,
		// which is composed of voters and articles waiting to vote
		err = openRandomAnonymousSpace(tx, _data.Comment, _version, _voterIDs)
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

func openRandomAnonymousSpace(tx *ent.Tx, comment string, vers *ent.Version, voterIDs []int) error {
	// 创建一个随机匿名空间
	_newRAS, err := tx.RAS.Create().
		SetComment(comment).
		SetVersion(vers).
		Save(ctx)
	if err != nil {
		return err
	}

	// 设置表决者初始投票状态：未投票
	_newVoters := make([]*ent.VoterCreate, len(voterIDs))
	for i, _voterID := range voterIDs {
		_newVoters[i] = tx.Voter.Create().SetRas(_newRAS).SetUserID(_voterID)
	}
	_, err = tx.Voter.CreateBulk(_newVoters...).Save(ctx)
	if err != nil {
		return err
	}

	// todo 倒计时，每个空间最长开启 72 小时，
	// 72 小时后，所有未投票的表决者，
	// 视为弃权，并触发警示系统。
	//
	// todo 2.0 争议性方案
	// 在某些领域的文章，
	// 是否可能需要花费远超 72 小时才能得到结论？
	// 72 小时作为固定时长，是否合理？

	return nil
}

func startSpaceCountdown() {

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

	_wordIDs := make([]int, len(vers.Edges.Keywords))
	for i, keyword := range vers.Edges.Keywords {
		_wordIDs[i] = keyword.ID
	}

	_userwordIDs, err := tx.UserWord.Query().
		Where(
			userword.HasWordWith(word.IDIn(_wordIDs...)),
			userword.HasUserWith(user.ID(userID))).
		Select(userword.FieldID).
		Ints(ctx)
	if err != nil {
		return err
	}

	_insertCount := len(_wordIDs) - len(_userwordIDs)

	if 0 < _insertCount {
		_wordBulk := make([]*ent.UserWordCreate, _insertCount)
		i := 0
		for _, keywordID := range _wordIDs {
			for _, userwordID := range _userwordIDs {
				if keywordID == userwordID {
					continue
				}
			}
			_wordBulk[i] = tx.UserWord.Create().SetUserID(userID).SetWordID(keywordID)
			i++
		}

		_, err = tx.UserWord.CreateBulk(_wordBulk...).Save(ctx)

		if err != nil {
			return err
		}
	}

	_, err = tx.UserWord.Update().Where(userword.IDIn(_userwordIDs...)).AddWeight(1).Save(ctx)

	if err != nil {
		return err
	}

	// _, err = tx.User.Update().
	// 	Where(user.ID(userID)).
	// 	AddWords(vers.Edges.Keywords...).
	// 	Save(ctx)

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
