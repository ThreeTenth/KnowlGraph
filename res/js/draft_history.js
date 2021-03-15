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
  },

  created() {
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id, needHistories: true }),
    }).then(function (resp) {
      _this.snapshots = resp.data.edges.Snapshots
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  beforeRouteEnter(to, from, next) {
    if (logined) {
      next(vm => {console.log(vm.snapshots)})
    } else {
      router.push({ name: "login" })
    }
  },

  template: fgm_draft_histories,
}

const DraftHistory = {
  props: ['id', 'hid', '__snapshot'],

  data: function() {
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
    }
  },

  created() {
    if (null != this.snapshot) {
      return
    }

    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/draft", { id: this.id, historyID: this.hid }),
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

  beforeRouteEnter(to, from, next) {
    if (logined) {
      next()
    } else {
      router.push({ name: "login" })
    }
  },

  template: fgm_draft_history,
}