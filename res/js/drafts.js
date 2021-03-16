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
        let snapshot = element.edges.Snapshots[0]
        let body = snapshot.body
        if (200 < body.length) {
          body = body.substr(0, 120).trim() + '...'
        }

        let createdat = this.i18n.CreatedAt1
        var options = { weekday: "long", year: "numeric", month: "short", day: "numeric", hour: "numeric", minute: "numeric", hour12: false };
        let created_at = new Date(snapshot.created_at).toLocaleString(getUserLang(), options)
        _drafts[index] = { body: body, created_at: createdat + created_at }
      }
      return _drafts
    }
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

  methods: {
    onDraft(i) {
      let draft = this.$data.__original[i]
      router.push({ name: 'editDraft', params: { id: draft.id, __draft: draft } })
    }
  },

  template: fgm_user_drafts,
}