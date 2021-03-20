// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: "private" }),
    }).then(function (resp) {
      router.push({ name: 'editDraft', params: { id: resp.data.id, __draft: resp.data } })
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
  },

  template: "",
}

const EditDraft = {
  props: ['id', "__draft"],

  data: function () {
    return {
      draft: this.__draft,
      __last: { body: "" },
    }
  },

  computed: {
    body: {
      get: function () {
        if (!this.draft) {
          return ""
        }
        if (this.draft.edges.Snapshots[0])
          return this.draft.edges.Snapshots[0].body
        return ""
      },
      set: function (newValue) {
        if (this.draft.edges.Snapshots[0]) {
          this.draft.edges.Snapshots[0].body = newValue
        } else {
          this.draft.edges.Snapshots[0] = { body: newValue }
        }
      }
    },
  },

  methods: {
    onHistories() {
      router.push({ name: 'draftHistories', params: { id: this.id } })
    },

    onPublish() {
      router.push({ name: 'publishArticle', params: { id: this.id, __draft: this.draft } })
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
      let content = this.draft.edges.Snapshots[0]
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
        _this.draft.edges.Snapshots[0] = resp.data
        _this.__setLast(_this.draft)
        // todo the new article have content and notify drafts page
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __load(__id, __draft) {
      if (__draft) {
        this.draft = __draft
        this.__setLast(__draft)
        return
      }

      if (this.draft && __id == this.draft.id) {
        this.__setLast(this.draft)
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: __id }),
      }).then(function (resp) {
        _this.draft = resp.data
        _this.__setLast(resp.data)
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __setLast(__draft) {
      if (__draft.edges.Snapshots[0]) {
        this.__last = __draft.edges.Snapshots[0].body
      } else {
        this.__last = ""
      }
    }
  },

  activated() {
    this.__load(this.id, this.__draft)
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id, to.params.__draft)
  },

  template: fgm_new_article,
}