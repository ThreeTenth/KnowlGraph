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
          id: article.id,
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
    onArticle(i) {
      let id = this.articles[i].id
      let title = this.articles[i].title
      if (!title) {
        title = this.articles[i].gist
      }
      title = title.trim().replace(/[;,/?:@&=+$_.!~*'()# ]+/g, '-')
      title = title.replace(/-$/g, '').toLowerCase()
      router.push({
        name: 'article', params: {
          id: id,
          title: title
        }
      })
    }
  },

  template: fgm_home,
}