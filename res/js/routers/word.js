// word.js

const Word = {
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
      this.qrcodeStatus = getStatus4Error(err)
      console.log(err);
    },
  },

  created() {
    getArticleByWord(this.id).then(this.loadData).catch(this.loadError)
  },

  template: fgm_word,
}