// index.js

const Index = {
  data: function () {
    return {
      pageStatus: 0,
      __original: [],
    }
  },

  computed: {
    articles: function () {
      let original = this.$data.__original
      if (!original) return []
      var _articles = []
      for (let index = 0; index < original.length; index++) {
        const version = original[index];
        const article = version.edges.article
        let title = version.title ? version.title : ""
        let code = title
        if (!code) {
          code = version.gist
        }
        code = encodeURLTitle(code)

        _articles[index] = {
          id: article.id,
          status: article.status,
          title: title,
          gist: version.gist,
          code: code,
          created_at: version.created_at
        }
      }
      return _articles
    }
  },

  created() {
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/articles", { limit: 10, offset: 0 }),
    }).then(function (resp) {
      _this.$data.__original = resp.data
      _this.pageStatus = resp.status
    }).catch(function (err) {
      _this.pageStatus = getStatus4Error(err)
    })
  },

  methods: {
    onArticle(i) {
      let id = this.articles[i].id
      let code = this.articles[i].code
      router.push({
        name: 'article', params: {
          id: id,
          code: code
        }
      })
    }
  },

  template: fgm_home,
}