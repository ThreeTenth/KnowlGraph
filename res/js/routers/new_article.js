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

      version: {
        title: "",
        gist: "",
        lang: getUserLang(),
        keywords: [],
        name: "v1.0.0",
        comment: "初次提交",
        init: false,
      },

      editingStatus: 0,
      fullscreen: false,

      isEditor: true,

      scroller: document,
      selectRange: {
        start: 0,
        end: 0,
      },
    }
  },

  watch: {
    id(value, oldVal) {
      console.log(value);
      this.reload()
    },
  },

  computed: {
    status: function () {
      if (!this.draft) { return "" }
      return this.draft.edges.article.status
    },

    preview: function () {
      return this.md2html(this.content)
    },

  },

  methods: {
    onHistories() {
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: this.id, needHistory: true }),
      }).then((resp) => {
        this.snapshots = resp.data.edges.snapshots
        this.snapshotIndex = 0
        this.showSnapshots = true
        this.__snapDiff(this.snapshotIndex, 1)
      }).catch((err) => {
        this.snapshotIndex = 0
        this.showSnapshots = true
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

    pasteEent(e) {
      const selection = window.getSelection()
      if (!selection.rangeCount) return
      if (!isChildAt(this.$refs.editor, selection.anchorNode)) return

      e.stopPropagation()
      e.preventDefault()

      let clipboardData = (e.clipboardData || window.clipboardData)
      let paste = ""
      let dataMode = null
      let codeLanguage = null

      for (let index = 0; index < clipboardData.items.length; index++) {
        const element = clipboardData.items[index]
        if (element.kind == "file") continue

        if (2 <= index) {
          paste = clipboardData.getData("text/plain").toString()
          dataMode = "text/plain"
          if (element.type == "vscode-editor-data") {
            const vscodeEditorData = clipboardData.getData(element.type)
            /* 
             * {"version":1,"isFromEmptySelection":false,"multicursorText":null,"mode":"javascript"} 
             */
            codeLanguage = JSON.parse(vscodeEditorData).mode
          }
        } else {
          paste = clipboardData.getData(element.type)
          dataMode = element.type
        }
      }

      if (codeLanguage) {
        paste = "```" + codeLanguage + "\n" + paste + "\n```\n"
      }
      switch (dataMode) {
        case "text/html":
          paste = this.turndownService.turndown(paste)
          break;

        default:
          paste = paste.replaceAll("\r\n", "\n")
          break;
      }
      document.execCommand("insertText", false, paste)
      // selection.deleteFromDocument()
      // selection.getRangeAt(0).insertNode(insertEl)
    },

    selectionchangeEvent(e) {
      let selection = document.getSelection()
      const editor = this.$refs.editor

      this.isEditor = isChildAt(editor, selection.focusNode)

      if (!this.isEditor) return

      var nodes = this.$refs.editor.childNodes
      var start = 0, end = 0
      for (let index = 0; index < nodes.length; index++) {
        const element = nodes[index];
        if (isChildAt(element, selection.anchorNode)) {
          start += selection.anchorOffset
          break
        } else {
          start += calcElementLength(element)
        }
      }
      for (let index = 0; index < nodes.length; index++) {
        const element = nodes[index];
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
    },

    onKeywordsChanged(values) {
      this.version.keywords = values
    },

    onShowPublish() {
      if (!this.content && this.version.init) {
        this.showPublish = true
        return
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

      var lang = this.version.lang
      var name = this.version.name
      var keywords = this.version.keywords
      var comment = this.version.comment
      var _original = this.draft.edges.original
      if (_original) {
        lang = _original.lang ? _original.lang : lang
        name = _original.name ? versionNamePlusOne(_original.name) : name
        keywords = []
        if (_original.edges.keywords) {
          _original.edges.keywords.forEach(element => {
            keywords.push(element.name)
          })
        }
        comment = ""
      }

      this.version.title = title
      this.version.gist = gist
      this.version.lang = lang
      this.version.keywords = keywords
      this.version.name = name
      this.version.comment = comment
      this.version.init = true
      this.showPublish = true
    },

    onPublish() {
      var draft = this.draft
      var version = this.version

      axios({
        method: "PUT",
        url: queryRestful("/v1/publish/article"),
        data: {
          name: version.name,
          comment: version.comment,
          title: version.title,
          gist: version.gist,
          lang: version.lang,
          keywords: version.keywords,
          draft_id: draft.id,
        },
      }).then((resp) => {
        // todo this content is published and notify home page and drafts page
        if (204 == resp.status) {
          router.push({ path: '/' })
          return
        }
        router.push({ name: 'achiveSelf' })
      }).catch(function (err) {
        console.log(err)
      })
    },

    onChanged(e) {
      this.content = textContentOfDiv(e.target)
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(() => { this.onSaveDraft() }, 2000)
    },

    onBlur() {
      window.clearTimeout(postChangedTimeoutID)
      this.content = textContentOfDiv(this.$refs.editor)
      this.onSaveDraft()
    },

    onSaveDraft: function () {
      // console.log(this.content);
      // console.log(this.__last);
      if (this.content === this.__last) {
        return
      }

      let temp = this._last
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/content"),
        data: {
          body: this.content,
          draft_id: this.draft.id,
        },
      }).then((resp) => {
        this.__setLast(resp.data.body)
        // todo the new article have content and notify drafts page
        // 编辑草稿时，所在草稿箱的草稿也会更新，不需要重新摘取请求了。
        // 
        // todo 保存状态
        //
        // todo 离线保存
      }).catch(() => {
        this.__setLast(temp)
      })
    },

    __setDraft(__draft) {
      this.draft = __draft
      this.content = getString(this.draft.edges.snapshots[0].body)
      setTimeout(() => {
        const editor = this.$refs.editor
        editor.focus()

        if ('' == this.content) {
          editor.innerText = ''
          return
        }
        var parts = this.content.split('\n')
        parts.forEach(part => {
          var el = document.createElement('div')
          el.innerText = "" === part ? "\n" : part
          editor.appendChild(el)
        });
        moveCursorToEnd(editor)
      }, 0);
    },

    __setLast(content) {
      this.__last = getString(content)
    },

    reset() {
      this.showSnapshots = false
      this.snapshots = []
      this.snapdiff = []
      this.snapshotIndex = 0

      this.showPreview = false
      this.showPublish = false

      this.content = ''
      this.draft = null
      this.__last = ''

      this.version = {
        title: "",
        gist: "",
        lang: getUserLang(),
        keywords: [],
        name: "v1.0.0",
        comment: "初次提交",
        init: false,
      }

      this.editingStatus = 0
      this.fullscreen = false

      this.isEditor = true

      this.scroller = document
      this.selectRange = {
        start: 0,
        end: 0,
      }
    },

    reload() {
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: this.id }),
      }).then((resp) => {
        this.reset()
        this.__setDraft(resp.data)
        this.__setLast(this.content)
        this.editingStatus = resp.status
      }).catch((err) => {
        console.error(err);
        this.editingStatus = getStatus4Error(err)
      })
    },
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
    this.reload()

    document.addEventListener("fullscreenchange", this.fullscreenEvent)
    document.addEventListener("paste", this.pasteEent)
    document.addEventListener('selectionchange', this.selectionchangeEvent);
  },

  beforeDestroy() {
    document.removeEventListener("fullscreenchange", this.fullscreenEvent)
    document.removeEventListener("paste", this.pasteEent)
    document.removeEventListener("selectionchange", this.selectionchangeEvent)
  },

  template: fgm_new_article,
}