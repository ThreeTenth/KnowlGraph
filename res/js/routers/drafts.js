// drafts.js

const Drafts = {
  data: function () {
    return {
      pageStatus: 0,
      __original: [],
    }
  },

  computed: {
    drafts: function () {
      let original = this.$data.__original
      var _drafts = []
      for (let index = 0; index < original.length; index++) {
        const element = original[index];
        let snapshot = element.edges.snapshots[0]
        let body = removeMarkdown(snapshot.body)
        let ps = body.split("\n")
        let firstP = ps[0]
        if (200 < firstP.length) {
          firstP = firstP.substr(0, 120).trim() + '...'
        }

        _drafts[index] = { id: element.id, body: firstP, created_at: snapshot.created_at }
      }
      return _drafts
    }
  },

  methods: {
    onDraft(i) {
      let draft = this.$data.__original[i]
      router.push({ name: 'editDraft', params: { id: draft.id, __draft: draft } })
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
      url: queryRestful("/v1/drafts"),
    }).then(function (resp) {
      _this.$data.__original = resp.data
      _this.pageStatus = resp.status
    }).catch(function (err) {
      _this.pageStatus = getStatus4Error(err)
    })
  },

  template: fgm_user_drafts,
}