// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: "private" }),
    }).then(function (resp) {
      router.push({ name: 'editDraft', params: { id: resp.data.DraftID, __new: true } })
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
  },

  template: "",
}

const EditDraft = {
  props: ['id', "__new"],

  data: function () {
    return {
      draft: {
        body: "",
      },
      __last: "",
    }
  },

  methods: {
    onHistories: function() {
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
      if (this.draft.body === this.__last) {
        return
      }

      console.log(this.draft.body)

      this.__last = this.draft.body
    },
  },

  created() {
    if (this.__new) {
      return
    }

    console.log(this.id)
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id }),
    }).then(function (resp) {
      let snapshots = resp.data.edges.Snapshots
      if (0 == snapshots.length) {
        return
      }
      _this.draft.body = snapshots[0].body
      _this.__last = snapshots[0].body
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  beforeRouteEnter(to, from, next) {
    if (logined) {
      next()
    } else {
      router.push({ name: "login" })
    }
  },

  template: fgm_new_article,
}