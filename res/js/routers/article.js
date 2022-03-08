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
      isNewNode: false,
      isNewQuote: false,
      isAddResponse: false,
      isSelectNode: false,
      selectNode: null,
      archiveMenuStatus: 0,
      pageStatus: 0,
      selection: {
        text: '',
        context: '',
        left: 0,
        top: 0,
      },
      nodes: null
    }
  },

  computed: {
    article: function () {
      let article = this.$data.__original
      if (undefined == article.edges)
        return {
          reactions: [],
          keywords: [],
          nodes: [],
          private: true,
        }

      let version = article.edges.versions[0]
      let content = version.edges.content

      var gist = version.gist
      var title = getTitleIfEmpty(version.title, version.gist)
      var code = encodeURLTitle(title)

      var body = content.body;
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
      var archives = article.edges.archives
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
        versionId: version.id,
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
        archives: archives,
      }
    }
  },

  methods: {
    onSelectNode(node) {
      this.isSelectNode = true
      this.selectNode = node
    },

    onPutNodeArticle() {
      putNodeArticle(this.selectNode.id, this.article.id)
        .then((resp) => {
          this.toast("归档成功", "success")
          this.selectNode.archive = true
          this.isSelectNode = false
          this.selectNode = null
          this.updateNodes()
        })
        .catch((err) => {
          this.toast(err, "error")
          this.isSelectNode = false
          this.selectNode = null
        })
    },

    onEditArticle() {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]

      axios({
        method: "PUT",
        url: queryRestful("/v1/article/edit", { id: version.id }),
      }).then(function (resp) {
        router.push({ name: 'editDraft', params: { id: resp.data.id } })
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

    onNewNode() {
      this.isNewNode = !this.isNewNode
    },

    newNodeSuccessed(archive) {
      this.isNewNode = false
      this.$data.__original.edges.archives.push(archive)
      this.selectNode = archive.edges.node
      this.onPutNodeArticle()
      this.addArchive(archive)
      this.updateNodes()
    },

    newNodeFailured(httpStatus, data) {
      this.isNewNode = false
      this.toast(data + "(" + httpStatus + ")", "error")
    },

    updateNodes() {
      var nodes = []
      var articleNodes = this.article.nodes || []
      var archives = this.article.archives || []
      articleNodes.forEach(node => {
        let archive = archives.find(item => {
          return item.edges.node.id == node.id
        })
        if (node.status == "private") {
          node.private = true
        } else {
          node.pub = true
        }
        if (archive) {
          node.archive = true
        }
        nodes.push(node)
      })
      archives = this.user.archives || []
      archives.forEach(archive => {
        let node = nodes.find(item => {
          return item.id == archive.edges.node.id
        })
        if (!node) {
          node = archive.edges.node
          nodes.push(node)
        }
        if (archive.status == "star") {
          node.star = true
        }
        if (archive.status == "watch") {
          node.watch = true
        }
      })

      if (nodes.length == 0) {
        return null
      }

      this.nodes = nodes
    },

    onStar() {
      if (this.user.logined) {
        if (this.article.star) {
          this.__deleteAsset('star')
        } else {
          this.__putAsset('star')
        }
      } else {
        router.push({ name: 'login' })
      }
    },

    onWatch() {
      if (this.user.logined) {
        if (this.article.watch) {
          this.__deleteAsset('watch')
        } else {
          this.__putAsset('watch')
        }
      } else {
        router.push({ name: 'login' })
      }
    },

    onAddQuote() {
      this.isAddResponse = true
      console.log("on add quote", this.selection);
    },

    switchLang() { },

    __putQuote() {
      let text = this.selection.text
      let context = this.selection.context
      let articleId = this.article.id
      let responseId = -1
      let highlight = 1
      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/quote"),
        data: {
          text: text,
          context: context,
          articleId: articleId,
          responseId: responseId,
          highlight: highlight,
        },
      }).then(function (resp) {
        _this.toast("Success", 'success')
      }).catch(function (err) {
        _this.toast("Failure", 'error')
      })
    },

    __resetSelection() {
      this.selection.text = ''
      this.selection.context = ''
    },

    __deleteAsset(status) {
      var assets = this.$data.__original.edges.assets
      var assetID = 0
      var index = 0
      for (; index < assets.length; index++) {
        const element = assets[index];
        if (element.status == status) {
          assetID = element.id
          break
        }
      }
      let _this = this
      axios({
        method: "DELETE",
        url: queryRestful("/v1/asset", { assetId: assetID }),
      }).then(function (resp) {
        assets.splice(index)
        _this.$vs.notification({
          color: 'success',
          position: 'bottom-right',
          title: "Success",
        })
      }).catch(function (resp) {
        _this.$vs.notification({
          color: 'danger',
          position: 'bottom-right',
          title: "Failure",
          text: resp.data
        })
      })
    },

    __putAsset(status) {
      let _this = this
      let params = { articleId: this.article.id, status: status }
      if (status == "browse") {
        params.versionId = this.article.versionId
      }
      axios({
        method: "PUT",
        url: queryRestful("/v1/asset", params),
      }).then(function (resp) {
        var assets = _this.$data.__original.edges.assets
        assets.push(resp.data)
        if (status == "browse") return
        _this.$vs.notification({
          color: 'success',
          position: 'bottom-right',
          title: "Success",
        })
      }).catch(function (resp) {
        if (status == "browse") return
        _this.$vs.notification({
          color: 'danger',
          position: 'bottom-right',
          title: "Failure",
          text: resp.data
        })
      })
    },

    __addQuote(e) {
      if (this.selection.text) {
        this.selection.left = e.pageX
        this.selection.top = e.pageY
        this.isNewQuote = true
      } else {
        this.__cancelQuote()
      }
    },

    __cancelQuote(e) {
      this.isNewQuote = false
      this.__resetSelection()
    },
  },

  created() {
    axios({
      method: "GET",
      url: queryRestful("/v1/article", { id: this.id }),
    }).then((resp) => {
      this.$data.__original = resp.data
      if (!this.code) {
        router.replace({
          name: 'article', params: {
            id: this.id,
            code: this.article.code,
          }
        })
      }
      if (this.user.logined) {
        this.$nextTick(() => {
          this.__putAsset("browse")
        })
      }
      this.pageStatus = resp.status
      this.updateNodes()
      document.title = this.article.title + " -- KnowlGraph"
    }).catch(function (err) {
      var status = getStatus4Error(err)
      if (401 == status && !logined) {
        router.push({ name: "login" })
        return
      }
      this.pageStatus = status
    })

    document.addEventListener('selectionchange', () => {
      let selection = document.getSelection()

      if (!this.$refs.reader
        || (selection.anchorNode == selection.focusNode && selection.anchorOffset == selection.focusOffset)) {
        this.__resetSelection()
        return
      }

      var text = selection.toString()
      if (!text || 0 == text.length || 0 == text.trim().length) return

      var divs = this.$refs.reader.children
      var startEl, endEl
      for (let index = 0; index < divs.length; index++) {
        const element = divs[index];
        if (isChildAt(element, selection.anchorNode)) {
          if (startEl) {
            endEl = element
          } else {
            startEl = element
          }
        }
        if (isChildAt(element, selection.focusNode)) {
          if (startEl) {
            endEl = element
          } else {
            startEl = element
          }
        }
      }

      var context = undefined
      for (let index = 0; index < divs.length; index++) {
        const element = divs[index];
        if (startEl == element) {
          context = element.innerText
          if (startEl == endEl) {
            break
          }
        } else if (endEl == element) {
          context += element.innerText
          break
        } else if (context) {
          context += element.innerText
        }

        context += "\n\n"
      }

      this.selection.text = selection.toString()
      this.selection.context = context
    })

    // document.onclick = this.__addQuote
  },

  template: fgm_article,
}