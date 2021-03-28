// publish_article.js

const PublishArticle = {
  props: ['id', '__draft'],

  data: function () {
    return {
      draft: {},
      title: "",
      gist: "",
      versionName: "",
      comment: "",
      lang: getUserLang(),
      keywords: [],
      words: [],
      languages: languages,
    }
  },

  watch: {
    draft: function (val, oldVal) {
      if (!val) return ""

      let content = val.edges.Snapshots[0].body
      let text = removeMarkdown(content)

      if (200 <= text.length) {
        this.gist = text.substr(0, 120).trim().replace(/[\r\n]/g, ' ') + '...'
      } else {
        this.gist = text
      }

      let found = content.match(/^# (.*)/m)

      if (found) {
        this.title = found[1]
        return
      }

      found = content.match(/^#+ (.*)/m)

      if (found) {
        this.title = found[1]
        return
      }
    },
  },

  methods: {
    onPublish() {
      var _this = this

      axios({
        method: "PUT",
        url: queryRestful("/v1/publish/article"),
        data: {
          name: this.versionName,
          comment: this.comment,
          title: this.title,
          gist: this.gist,
          lang: this.lang,
          keywords: this.keywords,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        // todo this content is published and notify home page and drafts page
        router.push({
          name: 'article', params: {
            id: resp.data.edges.Article.id,
            code: encodeURLTitle(_this.title)
          }
        })
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __load(__id, __draft) {
      if (__draft) {
        this.draft = __draft
        return
      }

      if (this.draft && __id == this.draft.id) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: __id }),
      }).then(function (resp) {
        _this.draft = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })

      axios({
        method: "GET",
        url: queryRestful("/v1/words"),
      }).then(function (resp) {
        _this.words = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  activated() {
    this.__load(this.id, this.__draft)
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id, to.params.__draft)
  },

  template: fgm_publish_article,
}