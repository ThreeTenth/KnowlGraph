// new_account.js

const NewAccount = {

  data: function () {
    return {
      requestState: 0,
      challenge: "",
      challengeState: 0,
      isQRCode: true,
    }
  },

  methods: {

    switchAuthMethod() {
      var is = this.isQRCode
      this.isQRCode = !is
    },

    onChallenge(challenge, requestState) {
      this.challenge = challenge
      this.requestState = requestState
    },

    onResult(qrcodeResult) {
      var state = qrcodeResult.state
      if (this.challengeState == state) {
        return
      }

      if (state == 4) {
        this.makeCredential(this.requestState, this.challenge)
      }

      this.challengeState = state
    },

    makeCredential(state, challenge) {
      axios({
        method: "PUT",
        url: queryRestful("/v1/account/t/makeCredential", { state: state, challenge: challenge })
      }).then(function (resp) {
        authSuccess(resp.data)
      }).catch(function (err) { })
    },

    createAccount() {
      let _this = this;
      let state = Math.random().toString(36).slice(2)
      axios({
        method: "GET",
        url: queryRestful("/v1/account/t/challenge", { state: state }),
      }).then(function (resp) {
        _this.makeCredential(state, resp.data)
      }).catch(function (err) {
        _this.toast(err, "error")
      })
    },
  },

  created() {
  },

  beforeDestroy() {
  },

  template: fgm_new_account,
}