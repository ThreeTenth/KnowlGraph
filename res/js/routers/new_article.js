// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: to.query.status }),
    }).then(function (resp) {
      router.push({ name: 'editDraft', params: { id: resp.data.id } })
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
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
      __last: { body: "" },
    }
  },

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
      }).catch(function (resp) {
        console.log(resp)
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
      var diff = Diff.diffChars(_old, _new)
      diff.forEach(part => {
        if (part.added || part.removed) {
          part.value = part.value.replace(/\n(?!.)/g, "\n ")
        }
      });
      this.diff = diff
    },

    onPublish() {
      this.showPreview = true
      // this.showPublish = true
      // router.push({ name: 'publishArticle', params: { id: this.id } })
    },

    onChanged: function (e) {
      this.content = e.target.value
      this.draft.edges.snapshots[0].body = e.target.value
      let _this = this
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(function () { _this.onSaveDraft() }, 2000)
    },

    onBlur: function () {
      window.clearTimeout(postChangedTimeoutID)
      this.onSaveDraft()
    },

    onSaveDraft: function () {
      if ("" === this.content) return
      if (this.__last && this.content && this.content === this.__last) {
        return
      }

      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/content"),
        data: {
          body: this.content,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        _this.draft.edges.snapshots[0] = resp.data
        _this.__setLast(_this.draft)
        // todo the new article have content and notify drafts page
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __setDraft(__draft) {
      this.draft = __draft
      if (!this.draft.edges.snapshots) {
        this.draft.edges.snapshots = [{ body: "" }]
      }
      this.content = this.draft.edges.snapshots[0].body
    },

    __setLast(__draft) {
      this.__last = __draft.edges.snapshots[0].body
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
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id }),
    }).then(function (resp) {
      _this.__setDraft(resp.data)
      _this.__setLast(resp.data)
    }).catch(function (resp) {
      console.log(resp)
    })
    document.body.style.overflow = 'hidden'
  },

  beforeDestroy() {
    document.body.style['overflow-y'] = 'scroll'
  },

  template: fgm_new_article,
}