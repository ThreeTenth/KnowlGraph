// publish_article.js

const PublishArticle = {
  data: function () {
    return {
      title: "",
      gist: "",
      versionName: "",
      comment: "",
      lang: getUserLang(),
      tags: [],
      languages: languages,
    }
  },

  methods: {
    onPublish() {
      console.log(this.$data)
    },
  },

  template: fgm_publish_article,
}