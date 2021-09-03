// login.js

const Login = {
  data: function () {
    return {
      syncStatus: 0,
      syncID: 0,
      challenge: 0,
    }
  },

  methods: {
    createAccount() {
      var _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/account/create"),
        data: {
          challenge: this.syncID,
        },
      }).then(function (resp) {
        Cookies.set("access_token", _this.syncID)
        window.open("/", "_self")
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
    let _this = this;
    axios({
      method: "GET",
      url: queryRestful("/v1/account/sync"),
    }).then(function (resp) {
      _this.syncID = resp.data
      _this.syncStatus = resp.status
      setTimeout(() => {
        new QRCode(_this.$refs.syncQrcode, "http://10.0.0.251:20010/" + _this.syncID);
        _this.$refs.syncQrcode.title = ''
      }, 0);
    }).catch(function (err) {
      _this.syncStatus = err.response.status
    })
  },

  template: fgm_login,
}