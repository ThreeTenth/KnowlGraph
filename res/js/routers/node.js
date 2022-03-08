// node.js

const NodeArticles = {
  props: ['id', 'name'],

  data: function () {
    return {
      pageStatus: 0,
      articles: [],
    }
  },

  methods: {
    loadData(response) {
      this.pageStatus = response.status
      this.articles = response.data
    },
    loadError(err) {
      this.pageStatus = getStatus4Error(err)
      console.error(err);
    },
  },

  created() {
    document.title = this.name + "下的所有文章 -- KnowlGraph"
    getNodeArticles(this.id).then(this.loadData).catch(this.loadError)
  },

  template: fgm_node
}