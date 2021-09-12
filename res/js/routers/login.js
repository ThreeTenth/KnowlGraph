// login.js

const Login = {
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
    createAccount() {
      if (this.user.logined) {
        this.toast("你已登录")
        return
      }

      var _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/account/create"),
        data: {
          challenge: this.syncID,
        },
      }).then(function (resp) {
        Cookies.set("access_token", _this.syncID)
        window.clearInterval(_this.checkChallengeInterval)
        window.open("/", "_self")
      }).catch(function (err) {
        window.clearInterval(_this.checkChallengeInterval)
        _this.toast("创建失败: " + err)
      })
    },

    checkChallengeState() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/check?challenge=" + this.syncID),
      }).then(function (resp) {
        var data = resp.data
        if (1 == data.state) {
          window.clearInterval(_this.checkChallengeInterval)
          _this.hasRefresh = true
        }
      }).catch(function (err) {
        window.clearInterval(_this.checkChallengeInterval)
        _this.hasRefresh = true
      })
    },

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
          let text = queryStatic("/t/" + resp.data)
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
        if (err.response) {
          _this.syncStatus = err.response.status
        } else {
          _this.syncStatus = 401
        }
      })
    },
  },

  created() {
    if (this.user.logined) {
      router.push({ path: '/' })
      return
    }
    this.getChallenge()
  },

  beforeDestroy() {
    window.clearInterval(this.checkChallengeInterval)
  },

  template: fgm_login,
}