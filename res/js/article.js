// article.js

const Article = {
  props: ['id', 'code'],

  data: function () {
    return {
      __original: {},
    }
  },

  computed: {
    article: function () {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.Versions[0]
      let content = version.edges.Content

      var title, gist, code = seo(version.title, version.gist)

      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      var body = converter.makeHtml(content.body);
      var keywords = version.edges.Keywords

      return {
        id: article.id,
        status: article.status,
        title: title,
        gist: gist,
        body: body,
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

      let version = article.edges.Versions[0]

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