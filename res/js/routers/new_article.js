// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: to.query.status }),
    }).then(function (resp) {
      router.push({ name: 'editDraft', params: { id: resp.data.id } })
    }).catch(function (err) {
      // console.log(resp.status, resp.data)
    })
  },

  template: "",
}

const EditDraft = {
  props: ['id'],

  data: function () {
    return {
      showSnapshots: false,
      snapshots: [],
      snapdiff: [],
      snapshotIndex: 0,

      showPreview: false,
      showPublish: false,

      content: '',
      draft: null,
      __last: '',

      editingStatus: 0,
      fullscreen: false,

      scroller: document,
      selectRange: {
        start: 0,
        end: 0,
      },
    }
  },

  // watch: {
  //   content(value, oldVal) {
  //     console.log(value);
  //   },
  // },

  computed: {
    status: function () {
      if (!this.draft) { return "" }
      return this.draft.edges.article.status
    },

    preview: function () {
      return this.md2html(this.content)
    },

    version: function () {
      if (!this.content) return {
        title: "",
        gist: "",
        lang: "",
        keywords: null,
        name: "",
        comment: "",
      }

      let title
      let gist

      let found = this.content.match(/^# (.*)/m)

      if (found) {
        title = found[1]
      } else {
        found = this.content.match(/^#+ (.*)/m)

        if (found) {
          title = found[1]
        }
      }

      var content = this.content.replace(/#+ .*/g, '')
      content = content.replace(/\n- .*/g, '')
      content = content.replace(/\n([123456789]+\.) .*/g, '')
      let text = removeMarkdown(content)

      if (!title) {
        var index = text.indexOf("\n")
        if (index === -1) index = text.length
        title = text.substring(0, index);
      }

      var index = text.indexOf(title)
      var gistStart = 0, gistEnd = 0

      if (index === -1) {
        gistStart = 0
        gistEnd = 120
      } else if (0 == index) {
        gistStart = title.length
        gistEnd = gistStart + 120
      } else {
        gistStart = 0
        gistEnd = index
      }

      if (text.length <= gistEnd) {
        gist = text.substr(gistStart, text.length).trim()
      } else {
        gist = text.substr(gistStart, gistEnd).trim() + '...'
      }

      index = gist.indexOf("\n")
      if (-1 < index) {
        gist = gist.substring(0, index)
      }

      if (!gist) {
        gist = title
      }

      var lang = getUserLang()
      var name = "v1.0.0"
      var keywords = null
      var comment = "初次提交："
      var _original = this.draft.edges.original
      if (_original) {
        lang = _original.lang
        name = _original.versionName ? _original.versionName : name
        keywords = []
        _original.edges.keywords.forEach(element => {
          keywords.push(element.name)
        });
        comment = ""
      }

      return {
        title: title,
        gist: gist,
        lang: lang,
        keywords: keywords,
        name: name,
        comment: comment,
      }
    },
  },

  methods: {
    onHistories() {
      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: this.id, needHistory: true }),
      }).then(function (resp) {
        _this.snapshots = resp.data.edges.snapshots
        _this.snapshotIndex = 0
        _this.showSnapshots = true
        _this.__snapDiff(_this.snapshotIndex, 1)
      }).catch(function (err) {
        _this.snapshotIndex = 0
        _this.showSnapshots = true
      })
    },

    onDafteHistory(i) {
      this.snapshotIndex = i
      this.__snapDiff(i, i + 1)
    },

    onHistoryPrev() {
      var i = this.snapshotIndex
      i--
      if (i < 0) {
        i = this.snapshots.length - 1
      }
      this.snapshotIndex = i
      this.__snapDiff(i, i + 1)
    },

    onHistoryNext() {
      var i = this.snapshotIndex
      i++
      if (this.snapshots.length <= i) {
        i = 0
      }
      this.snapshotIndex = i
      this.__snapDiff(i, i + 1)
    },

    __snapDiff(i, j) {
      var maxLen = this.snapshots.length
      var _new = i < maxLen ? this.snapshots[i].body : ''
      var _old = j < maxLen ? this.snapshots[j].body : ''
      var diff = Diff.diffChars(_old ? _old : '', _new ? _new : '')
      diff.forEach(part => {
        if (part.added || part.removed) {
          part.value = part.value.replace(/\n(?!.)/g, "\n ")
        }
      });
      this.diff = diff
    },

    insertEmoji(emoji) {
      document.execCommand("insertHTML", false, emoji)
    },

    onFullscreen() {
      if (this.fullscreen) {
        document.exitFullscreen()
      } else {
        this.$refs.editorContainer.requestFullscreen()
      }
      this.setFullscreen(!this.fullscreen)
    },

    fullscreenEvent(e) {
      if (document.fullscreenElement === null) {
        this.setFullscreen(false)
      }
    },

    setFullscreen(b) {
      this.fullscreen = b
      this.scroller = b ? this.$refs.editorContainer : document
    },

    onShowPublish() {
      this.showPreview = true
    },

    onPublish() {
      var _this = this

      axios({
        method: "PUT",
        url: queryRestful("/v1/publish/article"),
        data: {
          name: this.versionName,
          comment: this.comment,
          cover: this.cover,
          title: this.title,
          gist: this.gist,
          lang: this.lang,
          keywords: this.keywords,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        // todo this content is published and notify home page and drafts page
        if (204 == resp.status) {
          router.push({ path: '/' })
          return
        }
        router.push({
          name: 'article', params: {
            id: resp.data.edges.article.id,
            code: encodeURLTitle(_this.title)
          }
        })
      }).catch(function (err) {
        // console.log(err)
      })
    },

    onChanged(e) {
      var nodes = e.target.childNodes
      var content = ''
      for (let index = 0; index < nodes.length; index++) {
        const node = nodes[index]

        var text
        if (node.nodeName == "#text") {
          text = node.textContent
        } else {
          text = node.innerText
        }
        text = text.trim("\n")
        if (index + 1 < nodes.length) {
          text += '\n'
        }
        content += text
      }
      this.content = getString(content)
      // console.log(this.content);
      // this.content = e.target.innerText
      let _this = this
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(function () { _this.onSaveDraft() }, 2000)
    },

    onBlur() {
      window.clearTimeout(postChangedTimeoutID)
      this.onSaveDraft()
    },

    onSaveDraft: function () {
      // console.log(this.content);
      // console.log(this.__last);
      if (this.content === this.__last) {
        return
      }

      let temp = this._last
      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/content"),
        data: {
          body: this.content,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        _this.__setLast(resp.data.body)
        // todo the new article have content and notify drafts page
      }).catch(function (err) {
        _this.__setLast(temp)
      })
    },

    __setDraft(__draft) {
      this.draft = __draft
      this.content = getString(this.draft.edges.snapshots[0].body)
      if ('' == this.content) return
      var parts = this.content.split('\n')
      setTimeout(() => {
        parts.forEach(part => {
          var el = document.createElement('div')
          el.innerText = "" === part ? "\n" : part
          this.$refs.editor.appendChild(el)
        });
      }, 0);
    },

    __setLast(content) {
      this.__last = getString(content)
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
    document.title = "编辑 -- KnowlGraph"

    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id }),
    }).then(function (resp) {
      _this.__setDraft(resp.data)
      _this.__setLast(_this.content)
      _this.editingStatus = resp.status
    }).catch(function (err) {
      // console.log(err);
      _this.editingStatus = getStatus4Error(err)
    })

    document.addEventListener("fullscreenchange", this.fullscreenEvent)
    document.addEventListener("paste", (e) => {
      e.stopPropagation()
      e.preventDefault()
      var text = e.clipboardData.getData("Text")
      document.execCommand("insertText", false, text)
    })
    document.addEventListener('selectionchange', () => {
      let selection = document.getSelection()

      if (!_this.$refs.editor) return

      var divs = _this.$refs.editor.children
      var start = 0, end = 0
      for (let index = 0; index < divs.length; index++) {
        const element = divs[index];
        if (isChildAt(element, selection.anchorNode)) {
          start += selection.anchorOffset
          break
        } else {
          start += calcElementLength(element)
        }
      }
      for (let index = 0; index < divs.length; index++) {
        const element = divs[index];
        if (isChildAt(element, selection.focusNode)) {
          end += selection.focusOffset
          break
        } else {
          end += calcElementLength(element)
        }
      }
      if (start == end) {
        end = start + 1
      }
      // console.log(selection);
      // console.log(start, end);
      this.selectRange = {
        start: start,
        end: end,
      }
    });
  },

  beforeDestroy() {
    document.removeEventListener("fullscreenchange", this.fullscreenEvent)
  },

  template: fgm_new_article,
}