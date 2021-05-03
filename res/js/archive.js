// fgm_home.js

const Archive = {
  data: function () {
    return {
      status: "achiveSelf",
    }
  },

  created() {
    this.status = this.$router.history.current.name
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.status = this.$router.history.current.name
  },

  template: fgm_archive,
}

function getArchiveArticles(status) {
  return {
    data: function () {
      return {
        __original: [],
      }
    },

    computed: {
      articles: function () {
        let original = this.$data.__original
        if (original == undefined) return []
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
        url: queryRestful("/v1/assets", { status: status }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
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

    template: fgm_archive_articles,
  }
}

const archive_router = [
  { path: '', name: 'achiveSelf', component: getArchiveArticles("self"), props: true },
  { path: 'star', name: 'achiveStar', component: getArchiveArticles("star"), props: true },
  { path: 'watch', name: 'achiveWatch', component: getArchiveArticles("watch"), props: true },
  { path: 'readlist', name: 'achiveReadlist', component: getArchiveArticles("browse"), props: true },
]
