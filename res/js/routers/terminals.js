// terminals.js

const Terminals = {
  data: function () {
    return {
      pageStatus: 0,
      __original: [],
      syncID: "",
      showAuthorize: false,
      showConfirm: false,
      receiverName: null,
      confirmChallengeInterval: 0,
      confirmChallengeExpire: 15,
    }
  },

  computed: {
    terminals: function () {
      var parser = new UAParser();
      var orig = this.$data.__original
      var ts = []
      var tid = Cookies.get("terminal_id")
      for (let index = 0; index < orig.length; index++) {
        const term = orig[index];
        parser.setUA(term.ua);
        var result = parser.getResult();
        var device = result.device
        if (device.modal == undefined) {
          device = {
            model: "Desktop",
            type: "desktop",
            vendor: "Unknown",
          }
        }
        ts[index] = {
          current: term.id === parseInt(tid),
          nickname: term.name,
          name: result.browser.name,
          os: result.os.name,
          device: device,
        }
      }

      return ts
    },
  },

  methods: {
    os: function (prefix, terminal, suffix = "") {
      return prefix + terminal.os.toLowerCase() + suffix
    },

    dt: function (prefix, terminal, suffix = "") {
      return prefix + terminal.device.type.toLowerCase() + suffix
    },

    onAuthResult(syncID, data) {
      this.syncID = syncID
      this.receiverName = data.name
      this.showAuthorize = false
      this.showConfirm = true
      this.confirmChallengeInterval = window.setInterval(() => {
        if (0 < this.confirmChallengeExpire) {
          this.confirmChallengeExpire--
        } else {
          window.clearInterval(this.confirmChallengeInterval)
          this.showConfirm = false
        }
      }, 1000);
    },

    onToggleModal(toggle) {
      if (toggle) return
      this.onCancelAuthn("DELETE")
    },

    onCancelAuthn() {
      this.__postAccountAuthn("DELETE")
    },

    onTemporaryAuthn() {
      this.showConfirm = true
      this.__postAccountAuthn("PATCH")
    },

    onConfirmAuthn() {
      this.showConfirm = true
      this.__postAccountAuthn("POST")
    },

    __postAccountAuthn(method) {
      if ("" === this.syncID) return

      var _this = this
      axios({
        method: method,
        url: queryRestful("/v1/account/authn"),
        data: {
          challenge: this.syncID,
        },
      }).then(function (resp) {
        if ("DELETE" == method) return
        _this.toast("Success", "success")
      }).catch(function (err) {
        _this.toast("错误: " + err, "error")
      })
    },

    __updateTerminals() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/terminals"),
      }).then(function (resp) {
        _this.pageStatus = resp.status
        _this.$data.__original = resp.data
      }).catch(function (err) {
        _this.pageStatus = getStatus4Error(err)
      })
    },
  },

  created() {
    this.__updateTerminals()
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  beforeDestroy() {
    if (0 < this.confirmChallengeExpire) {
      window.clearInterval(this.confirmChallengeInterval)
    }
  },

  template: fgm_terminals,
}