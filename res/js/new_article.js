// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: "private" }),
    }).then(function (resp) {
      console.log(resp.status, resp.data)
      // articleID = resp.data.ArticleID
      // draftID = resp.data.DraftID
      // this.setContent()
      // next()
      // 命名的路由
      router.push({ name: 'editDraft', params: { id: resp.data.DraftID } })
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
  },

  template: "",
}

const EditDraft = {
  data: function () {
    return {
      draft: {
        body: ""
      }
    }
  },
  methods: {
    onChanged: function () {
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(this._postArticleContent(), 800)
    },

    onBlur: function () {
      window.clearTimeout(postChangedTimeoutID)
      this._postArticleContent()
    },

    _postArticleContent: function () {
      console.log(this.draft.body)
    },
  },
  template: fgm_new_article,
}