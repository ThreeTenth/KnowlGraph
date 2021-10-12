// vue_auth_qrcode.js

Vue.component('auth-qrcode', {

  data: function () {
    return {
      qrcodeStatus: 0,
      challenge: "",
      requestState: "",
      challengeState: 1,
      checkChallengeInterval: 0,
      hasRefresh: false,
      qrcode: null,
    }
  },

  methods: {
    getChallenge() {
      let _this = this;
      let state = Math.random().toString(36).slice(2)
      axios({
        method: "GET",
        url: queryRestful("/v1/account/t/challenge", { state: state }),
      }).then(function (resp) {
        _this.requestState = state
        _this.challenge = resp.data
        _this.qrcodeStatus = resp.status
        _this.hasRefresh = false
        _this.$emit('challenge', resp.data, state)
        setTimeout(() => {
          let text = queryPage("/g/" + resp.data)
          if (_this.qrcode) {
            _this.qrcode.clear()
            _this.qrcode.makeCode(text)
          } else {
            _this.qrcode = new QRCode(_this.$refs.qrcode, text)
          }
          _this.$refs.qrcode.title = ''

          _this.checkChallengeInterval = window.setInterval(() => {
            _this.checkChallengeState()
          }, 3000);
        }, 0);
      }).catch(function (err) {
        _this.qrcodeStatus = getStatus4Error(err)
      })
    },

    checkChallengeState() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/t/scanChallenge", {
          challenge: this.challenge
        }),
      }).then(function (resp) {
        var state = resp.data.state
        _this.challengeState = state
        if (4 == state) {
          window.clearInterval(_this.checkChallengeInterval)
          _this.hasRefresh = true
        }
        _this.$emit('result', state)
      }).catch(function (err) {
        window.clearInterval(_this.checkChallengeInterval)
        _this.hasRefresh = true
        _this.challengeState = 1
      })
    },

    onCancelAuthn() {
      this.hasRefresh = true
      axios({
        method: "DELETE",
        url: queryRestful("/v1/account/t/cancel", {
          challenge: this.challenge,
        }),
      })
    },
  },

  created() {
    this.getChallenge()
  },

  beforeDestroy() {
    window.clearInterval(this.checkChallengeInterval)
    if (this.hasRefresh || this.qrcodeStatus != 200) return

    this.onCancelAuthn()
  },

  template: com_auth_qrcode,
})