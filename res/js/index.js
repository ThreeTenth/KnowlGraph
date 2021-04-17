// fgm_home.js

const Index = {
  data: function () {
    return {
      __original: [],
      __published: [],
    }
  },

  computed: {
    articles: function () {
      let original = this.$data.__original.concat(this.$data.__published)
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
      url: queryRestful("/v1/user/articles", { status: "self" }),
    }).then(function (resp) {
      _this.$data.__original = resp.data
    }).catch(function (resp) {
      console.log(resp)
    })
    axios({
      method: "GET",
      url: queryRestful("/v1/articles", { limit: 10, offset: 0 }),
    }).then(function (resp) {
      _this.$data.__published = resp.data
    }).catch(function (resp) {
      console.log(resp)
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