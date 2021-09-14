// terminals.js

const Terminals = {
  data: function () {
    return {
      pageStatus: 0,
      terminals: [],
    }
  },

  methods: {
    checkChallengeState() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/check?challenge=" + this.syncID),
      }).then(function (resp) {
      }).catch(function (err) {
      })
    },

    getChallenge() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/sync"),
      }).then(function (resp) {

      }).catch(function (err) {
      })
    },

    onToggleModal(toggle) {
      if (toggle) return
      this.onCancelAuthn()
    },

    onCancelAuthn() {
      this.canAuthn = false
      this.hasRefresh = true
      this.__postAccountAuthn("DELETE")
    },

    onTemporaryAuthn() {
      this.canAuthn = false
      this.hasRefresh = true
      this.__postAccountAuthn("PATCH")
    },

    onConfirmAuthn() {
      this.canAuthn = false
      this.hasRefresh = true
      this.__postAccountAuthn("POST")
    },

    __postAccountAuthn(method) {
      var _this = this
      axios({
        method: method,
        url: queryRestful("/v1/account/authn"),
        data: {
          challenge: this.syncID,
        },
      }).then(function (resp) {
      }).catch(function (err) {
      })
    },

    __updateTerminals() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/terminals"),
      }).then(function (resp) {
        _this.pageStatus = resp.status
      }).catch(function (err) {
        _this.pageStatus = getStatus4Error(err)
      })
    },
  },

  created() {
    this.__updateTerminals()
  },

  beforeDestroy() {
  },

  template: fgm_terminals,
}