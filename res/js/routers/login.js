// login.js

const Login = {
  data: function () {
    return {
    }
  },

  methods: {
    onAuthResult(syncID, data) {
      this.todo("扫描成功，等待授权中")
    },

    createAccount(syncID) {
      var _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/account/create"),
        data: {
          challenge: syncID,
        },
      }).then(function (resp) {
        // Default(no expiration time is set): Cookie is removed when the user closes the browser.
        Cookies.set("terminal_id", resp.data, { expires: 365 })
        // Token 有效期最多 30 天。
        Cookies.set("access_token", syncID, { expires: 30 })
        window.open("/", "_self")
      }).catch(function (err) {
        _this.toast("创建失败: " + err)
      })
    },

    startCreate() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/sync"),
      }).then(function (resp) {
        _this.createAccount(resp.data)
      }).catch(function (err) {
        _this.toast("创建失败: " + err)
      })
    },
  },

  created() {
    if (this.user.logined) {
      router.push({ path: '/' })
      return
    }
  },

  beforeDestroy() {
  },

  template: fgm_login,
}