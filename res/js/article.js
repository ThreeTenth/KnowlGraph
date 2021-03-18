// article.js

const Article = {
  props: ['id', 'title'],

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

      return {
        id: article.id,
        status: article.status,
        title: version.title ? version.title : "",
        gist: version.gist,
        body: content.body,
        created_at: version.created_at
      }
    }
  },

  methods: {
    __load(id) {
      if (null != this.snapshot) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/article", { id: id }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.snapshot = to.params.__snapshot
      vm.__load(to.params.id)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.snapshot = to.params.__snapshot
    this.__load(to.params.id)
  },

  template: fgm_article,
}