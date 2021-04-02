// article.js

const Article = {
  props: ['id', 'code'],

  data: function () {
    return {
      emoji: {
        "up": "👍", "down": "👎", "laugh": "😄", "hooray": "🎉", "confused": "😕", "heart": "❤️", "rocket": "🚀", "eyes": "👀",
      },
      __original: {},
    }
  },

  computed: {
    article: function () {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]
      let content = version.edges.content

      var title, gist, code = seo(version.title, version.gist)

      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      var body = converter.makeHtml(content.body);
      var keywords = version.edges.keywords

      var lang = userLang
      languages.forEach(element => {
        if (element.code == version.lang) {
          lang = element
        }
      });

      var reactions = [
        { status: "up", count: "9" },
        { status: "down", count: "99" },
        { status: "laugh", count: "999" },
        { status: "heart", count: "1K" },
        { status: "hooray", count: "10K" },
        { status: "confused", count: "100k" },
        { status: "rocket", count: "1M" },
        { status: "eyes", count: "10M" },
      ]
      reactions = article.edges.reactions

      return {
        id: article.id,
        status: article.status,
        title: title,
        gist: gist,
        body: body,
        lang: lang,
        reactions: reactions,
        keywords: keywords,
        created_at: version.created_at,
        code: code,
      }
    }
  },

  methods: {
    onEditArticle() {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]

      axios({
        method: "PUT",
        url: queryRestful("/v1/article/edit", { id: version.id }),
      }).then(function (resp) {
        router.push({ name: 'editDraft', params: { id: resp.data.id, __draft: resp.data } })
      }).catch(function (resp) {
        console.log(resp.status, resp.data)
      })
    },

    __load(id) {
      if (this.$data.__original.id == this.id) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/article", { id: id }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
        if (!_this.code) {
          router.replace({
            name: 'article', params: {
              id: _this.id,
              code: _this.article.code,
            }
          })
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.__load(to.params.id)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id)
  },

  template: fgm_article,
}