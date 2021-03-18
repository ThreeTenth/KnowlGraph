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
        const element = original[index];
        let article = element.edges.Article
        let version = article.edges.Versions[0]

        _articles[index] = {
          status: article.status,
          title: version.title ? version.title : "",
          gist: version.gist,
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
  },

  methods: {
    onArticles(i) {
      let article = this.$data.__original[i]
      console.log(article)
      // router.push({ name: 'article', params: { id: article.id, __article: article } })
    }
  },

  template: fgm_home,
}