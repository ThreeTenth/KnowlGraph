// publish_article.js

const PublishArticle = {
  data: function() {
    return {
      user: {
        lang: getUserLang(),
      },
      languages: languages,
    }
  },

  methods: {
    onSelectArticleLang() {

    },
  },

  template: fgm_publish_article,
}