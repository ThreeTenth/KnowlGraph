// draft_history.js

const DraftHistory = {
  props: ['id', 'hid', '__snapshot'],

  data: function () {
    return {
      snapshot: this.__snapshot,
    }
  },

  methods: {
    onHistory: function () {
    },
    onBack: function () {
      if (3 <= window.history.length) {
        router.go(-1)
      } else {
        router.push({ name: 'draftHistories', params: { id: this.id } })
      }
    },
    __load(id, hid) {
      if (null != this.snapshot) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: id, historyID: hid }),
      }).then(function (resp) {
        let snapshots = resp.data.edges.snapshots
        if (snapshots) {
          _this.snapshot = snapshots[0]
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next(vm => {
      vm.snapshot = to.params.__snapshot
      vm.__load(to.params.id, to.params.hid)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.snapshot = to.params.__snapshot
    this.__load(to.params.id, to.params.hid)
  },

  template: fgm_draft_history,
}