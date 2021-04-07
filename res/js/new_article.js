// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: to.query.status }),
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

    __load(__id, __draft) {
      if (__draft) {
        this.__setDraft(__draft)
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
        _this.__setDraft(resp.data)
        _this.__setLast(resp.data)
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