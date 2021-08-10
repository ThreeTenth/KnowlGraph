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
      draft: null,
      __last: { body: "" },
    }
  },

  computed: {
    body: {
      get: function () {
        if (this.draft && this.draft.edges.snapshots)
          return this.draft.edges.snapshots[0].body
        return ""
      },
      set: function (newValue) {
        this.draft.edges.snapshots[0].body = newValue
      }
    },

    status: function () {
      if (!this.draft) { return "" }
      return this.draft.edges.article.status
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
      router.push({ name: 'publishArticle', params: { id: this.id } })
    },

    onChanged: function () {
      let _this = this
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(function () { _this.onSaveDraft() }, 2000)
    },

    onBlur: function () {
      window.clearTimeout(postChangedTimeoutID)
      this.onSaveDraft()
    },

    onSaveDraft: function () {
      // console.trace()
      let content = this.draft.edges.snapshots[0]
      if ("" === content.body) return
      if (this.__last && content.body && content.body === this.__last) {
        return
      }

      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/content"),
        data: {
          body: content.body,
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
  },

  template: fgm_new_article,
}