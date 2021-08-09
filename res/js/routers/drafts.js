// drafts.js

const Drafts = {
  data: function () {
    return {
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
        if (200 < body.length) {
          body = body.substr(0, 120).trim() + '...'
        }

        _drafts[index] = { id: element.id, body: body, created_at: snapshot.created_at }
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
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  template: fgm_user_drafts,
}