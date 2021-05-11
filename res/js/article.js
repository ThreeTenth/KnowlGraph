// article.js

const Article = {
  props: ['id', 'code'],

  data: function () {
    return {
      emoji: {
        "up": "👍", "down": "👎", "laugh": "😄", "hooray": "🎉", "confused": "😕", "heart": "❤️", "rocket": "🚀", "eyes": "👀",
      },
      reactions: [
        { "name": "Up", "value": "up", "emoji": "👍" },
        { "name": "Down", "value": "down", "emoji": "👎" },
        { "name": "Laugh", "value": "laugh", "emoji": "😄" },
        { "name": "Hooray", "value": "hooray", "emoji": "🎉" },
        { "name": "Confused", "value": "confused", "emoji": "😕" },
        { "name": "Heart", "value": "heart", "emoji": "❤️" },
        { "name": "Rocket", "value": "rocket", "emoji": "🚀" },
        { "name": "Eyes", "value": "eyes", "emoji": "👀" },
      ],
      __original: {},
    }
  },

  computed: {
    article: function () {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]
      let content = version.edges.content

      var title, gist, code = seo(version.title, version.gist)

      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      var body = converter.makeHtml(content.body);
      var keywords = version.edges.keywords

      var lang = userLang
      languages.forEach(element => {
        if (element.code == version.lang) {
          lang = element
        }
      });

      var reactions = article.edges.reactions
      var assets = article.edges.assets
      var nodes = article.edges.nodes
      var star = false
      var watch = false
      var private = false
      var browseCount = 0
      if (assets) {
        for (let index = 0; index < assets.length; index++) {
          const element = assets[index];
          if (element.status == "star") {
            star = true
          } else if (element.status == "watch") {
            watch = true
          } else if (element.status == "self") {
            private = true
          } else if (element.status == "browse") {
            browseCount++
          }
        }
      }

      // sort by value
      reactions.sort(function (a, b) {
        return (b.count - a.count)
      });

      return {
        id: article.id,
        status: article.status,
        title: title,
        gist: gist,
        body: body,
        lang: lang,
        reactions: reactions,
        keywords: keywords,
        created_at: version.created_at,
        code: code,
        star: star,
        watch: watch,
        private: private,
        browseCount: browseCount,
        nodes: nodes,
      }
    }
  },

  methods: {
    onEditArticle() {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]

      axios({
        method: "PUT",
        url: queryRestful("/v1/article/edit", { id: version.id }),
      }).then(function (resp) {
        router.push({ name: 'editDraft', params: { id: resp.data.id, __draft: resp.data } })
      }).catch(function (resp) {
        console.log(resp.status, resp.data)
      })
    },

    onPickReaction(reaction) {
      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/reaction", { articleId: this.article.id, reaction: reaction }),
      }).then(function (resp) {
        let reactions = _this.article.reactions
        let ok = false
        reactions = reactions == undefined ? [] : reactions
        for (let index = 0; index < reactions.length; index++) {
          const reac = reactions[index];
          if (reac.status == reaction) {
            reac.count += 1
            ok = true
            break
          }
        }

        if (!ok) {
          reactions.push(resp.data)
        }

        _this.$data.__original.edges.reactions = reactions
      }).catch(function (resp) {
        console.error(resp.status, resp.data)
      })
    },

    __load(id) {
      if (this.$data.__original.id == this.id) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/article", { id: id }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
        if (!_this.code) {
          router.replace({
            name: 'article', params: {
              id: _this.id,
              code: _this.article.code,
            }
          })
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.__load(to.params.id)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id)
  },

  template: fgm_article,
}