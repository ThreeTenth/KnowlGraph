// fgm_home.js

const Index = {
  data: function () {
    return {
      __original: [],
    }
  },

  computed: {
    articles: function () {
      let original = this.$data.__original
      var _articles = []
      for (let index = 0; index < original.length; index++) {
        const article = original[index];
        let version = article.edges.versions[0]
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
      url: queryRestful("/v1/articles"),
    }).then(function (resp) {
      console.log(resp.data)
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