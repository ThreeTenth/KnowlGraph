// vue_auth_qrcode.js

Vue.component('auth-qrcode', {

  data: function () {
    return {
      syncStatus: 0,
      syncID: 0,
      checkChallengeInterval: 0,
      hasRefresh: false,
      qrcode: null,
    }
  },

  methods: {
    getChallenge() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/sync"),
      }).then(function (resp) {
        _this.syncID = resp.data
        _this.syncStatus = resp.status
        _this.hasRefresh = false
        setTimeout(() => {
          let text = queryPage("/g/" + resp.data)
          if (_this.qrcode) {
            _this.qrcode.clear()
            _this.qrcode.makeCode(text)
          } else {
            _this.qrcode = new QRCode(_this.$refs.syncQrcode, text)
          }
          _this.$refs.syncQrcode.title = ''

          _this.checkChallengeInterval = window.setInterval(() => {
            _this.checkChallengeState()
          }, 3000);
        }, 0);
      }).catch(function (err) {
        _this.syncStatus = getStatus4Error(err)
      })
    },

    checkChallengeState() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/check?challenge=" + this.syncID),
      }).then(function (resp) {
        var data = resp.data
        if (2 == data.state) {
          window.clearInterval(_this.checkChallengeInterval)
          _this.hasRefresh = true

          _this.$emit('result', _this.syncID, data)
        }
      }).catch(function (err) {
        window.clearInterval(_this.checkChallengeInterval)
        _this.hasRefresh = true
      })
    },

    onCancelAuthn() {
      this.hasRefresh = true
      axios({
        method: "DELETE",
        url: queryRestful("/v1/account/authn"),
        data: {
          challenge: this.syncID,
        },
      })
    },
  },

  created() {
    this.getChallenge()
  },

  beforeDestroy() {
    window.clearInterval(this.checkChallengeInterval)
    if (this.hasRefresh || this.syncStatus != 200) return

    this.onCancelAuthn()
  },

  template: com_auth_qrcode,
})