// vote.js

const Vote = {
  data: function() {
    return {
      showPromiseDialog: true,
    }
  },
  computed: {
    version: function() {
      var ver = this.vote.value.edges.version
      this.languages.forEach(element => {
        if (element.code == ver.lang) {
          ver.lang = element.name
        }
      });
      return ver
    },
  },
  methods: {
    onPromise() {
      this.showPromiseDialog = false
    },
    onAllow() {
      if (!this.vote.ras) return

      this.__postVote("allowed")
    },
    onRejecte() {
      if (!this.vote.ras) return

      this.__postVote("rejected")
    },
    __postVote(status) {
      var _this = this
      axios({
        method: "POST",
        url: queryRestful("/v1/vote"),
        data: {
          id: this.vote.ras.id,
          status: status,
        },
      }).then(function (resp) {
        _this.vote.ras = null
        _this.vote.has = false
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },
  created() { },
  template: fgm_vote,
}