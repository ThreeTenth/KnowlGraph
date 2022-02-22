// index.js

const Index = {
  data: function () {
    return {
      pageStatus: 0,
      articles: [],
    }
  },

  computed: {},

  created() {
    document.title = "Knowledge graph -- KnowlGraph"

    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/articles", { limit: 10, offset: 0 }),
    }).then(function (resp) {
      _this.$data.articles = resp.data
      _this.pageStatus = resp.status
    }).catch(function (err) {
      _this.pageStatus = getStatus4Error(err)
    })
  },

  methods: {},

  template: fgm_home,
}