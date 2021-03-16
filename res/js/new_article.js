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
      __last: "",
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
    onHistories: function () {
      router.push({ name: 'draftHistories', params: { id: this.id } })
    },

    onChanged: function () {
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(this.__postArticleContent(), 800)
    },

    onBlur: function () {
      window.clearTimeout(postChangedTimeoutID)
      this.__postArticleContent()
    },

    __postArticleContent: function () {
      let content = this.draft.edges.Snapshots[0].body
      if (content.body === this.__last) {
        return
      }

      this.__last = content.body
    },

    __load(id) {
      if (this.draft) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: id }),
      }).then(function (resp) {
        _this.draft = resp.data
        _this.__last = resp.data.edges.Snapshots[0].body
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    if (logined) {
      next(vm => {
        vm.draft = to.params.__draft
        vm.__load(to.params.id)
      })
    } else {
      router.push({ name: "login" })
    }
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.draft = to.params.__draft
    this.__load(to.params.id)
  },

  template: fgm_new_article,
}