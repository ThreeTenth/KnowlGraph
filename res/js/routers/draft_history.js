// draft_history.js

const DraftHistory = {
  props: ['id', 'hid'],

  data: function () {
    return {
      snapshot: null,
    }
  },

  created() {
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id, historyID: this.hid }),
    }).then(function (resp) {
      let snapshots = resp.data.edges.snapshots
      if (snapshots) {
        _this.snapshot = snapshots[0]
      }
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  methods: {
    onBack: function () {
      if (3 <= window.history.length) {
        router.go(-1)
      } else {
        router.push({ name: 'draftHistories', params: { id: this.id } })
      }
    },
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  template: fgm_draft_history,
}