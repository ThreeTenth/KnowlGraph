// login.js

const Login = {
  data: function () {
    return {
      syncStatus: 0,
      syncID: 0,
    }
  },

  methods: {
    onAuthResult(syncID, data) {
      this.todo("扫描成功，等待授权中")
    },

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

    getChallenge() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/sync"),
      }).then(function (resp) {
        _this.syncID = resp.data
        _this.syncStatus = resp.status
      }).catch(function (err) {
        _this.syncStatus = getStatus4Error(err)
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
  },

  template: fgm_login,
}