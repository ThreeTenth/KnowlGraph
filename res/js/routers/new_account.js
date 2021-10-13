// new_account.js

const NewAccount = {

  data: function () {
    return {
      requestState: 0,
      challenge: "",
      challengeState: 0,
    }
  },

  methods: {

    onChallenge(challenge, requestState) {
      this.challenge = challenge
      this.requestState = requestState
    },

    onResult(state) {
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
        const data = resp.data
        // Default(no expiration time is set): Cookie is removed when the user closes the browser.
        const idExpiresDay = data.onlyOnce ? 0 : 365
        const tokenExpiresDay = data.onlyOnce ? 0 : 30

        Cookies.set("terminal_id", data.id, { expires: idExpiresDay })
        // Token 有效期最多 30 天。
        Cookies.set("access_token", data.token, { expires: tokenExpiresDay })
        window.open("/", "_self")
      }).catch(function (err) {})
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
    axios({
      method: "POST",
      url: queryRestful("/v1/account/t/cancel", { challenge: this.challenge }),
    })
  },

  template: fgm_new_account,
}