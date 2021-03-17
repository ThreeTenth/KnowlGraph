// draft_history.js

const DraftHistories = {
  props: ['id'],

  data: function () {
    return {
      snapshots: [],
    }
  },

  methods: {
    onHistory: function (snapshot) {
      router.push({ name: 'draftHistory', params: { id: this.id, hid: snapshot.id, __snapshot: snapshot } })
    },

    __load(id) {
      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: id, needHistory: true }),
      }).then(function (resp) {
        _this.snapshots = resp.data.edges.Snapshots
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

    if (from.name === "draftHistory") {
      next()
    } else {
      next(vm => {
        vm.__load(to.params.id)
      })
    }
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id)
  },

  template: fgm_draft_histories,
}

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
        let snapshots = resp.data.edges.Snapshots
        if (0 == snapshots.length) {
          return
        }
        _this.snapshot = snapshots[0]
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