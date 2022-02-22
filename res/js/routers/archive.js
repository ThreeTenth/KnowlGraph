// archive.js

const Archive = {
  data: function () {
    return {
      status: "achiveSelf",
    }
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
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
        articles: [],
        status: status,
        icons: {
          self: "🔒",
          star: "⭐",
          watch: "👀",
          browse: "🕒",
        },
        pageStatus: 0,
      }
    },

    computed: {},

    created() {
      document.title = "归档 -- KnowlGraph"

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/assets", { status: status }),
      }).then(function (resp) {
        _this.$data.articles = resp.data
        _this.pageStatus = resp.status
      }).catch(function (err) {
        _this.pageStatus = getStatus4Error(err)
      })
    },

    methods: {
      onDelete() {

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
