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

      isCamera: true,
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
        if (device.model == undefined) {
          device = {
            model: "Desktop",
            type: "desktop",
            vendor: "Unknown",
          }
        }
        let current = term.id === parseInt(tid)
        ts[index] = {
          id: term.id,
          canDelete: false,
          current: current,
          nickname: term.name,
          name: result.browser.name,
          os: result.os.name,
          device: device,
          onlyOnce: term.onlyOnce,
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

    onAddTerminal() {
      var tid = Cookies.get("terminal_id")
      putBeginValidate(tid, (resp) => {
        this.beginValidateSuccess(tid, resp.data)
      }, this.beginValidateFailure)
    },

    beginValidateFailure(err) {
      this.toast(err, "error")
      console.info("BeginLoginFailure", err)
    },
    beginValidateSuccess(terminalID, makeAssertionOptions) {
      console.log("Assertion Options:");
      console.log(makeAssertionOptions);
      navigator.credentials.get(makeAssertionOptions)
        .then((credential) => {
          console.log(credential);
          postFinishValidate(terminalID, credential, this.finishValidateSuccess, this.finishValidateFailure)
        }).catch((err) => {
          this.toast(err, "error")
          console.info("Error", err)
        });
    },

    finishValidateFailure(err) {
      this.toast(err, "error")
      console.info("FinishLoginFailure", err)
    },
    finishValidateSuccess(resp) {
      console.log(resp.data)
      this.showAuthorize = true
    },

    canDel(terminal) {
      return !terminal.current || 1 == this.terminals.length
    },

    switchAuthMethod() {
      var isCamera = this.isCamera

      this.isCamera = !isCamera
    },

    onChallenge(challenge, requestState) {
      this.syncID = challenge
    },

    onAuthorizeTerminal(resp, challenge) {
      this.syncID = challenge
      this.onAuthResult(resp.data)
    },

    onAuthResult(qrcodeResult) {
      if (qrcodeResult.state != 2) return
      this.receiverName = qrcodeResult.name
      this.confirmChallengeExpire = 15
      this.showAuthorize = false
      this.showConfirm = true

      var confirmChallengeInterval = window.setInterval(() => {
        if (0 < this.confirmChallengeExpire) {
          this.confirmChallengeExpire--
        } else {
          window.clearInterval(confirmChallengeInterval)
          this.showConfirm = false
        }
      }, 1000);
      this.confirmChallengeInterval = confirmChallengeInterval
    },

    onDeleteTerminal(terminal) {
      var _this = this
      axios({
        method: "DELETE",
        url: queryRestful("/v1/account/terminal", { id: terminal.id }),
      }).then(function (resp) {
        if (terminal.current) {
          Cookies.remove("terminal_id")
          Cookies.remove("access_token")
          window.open("/", "_self")
        } else {
          _this.toast("Delete success", "success")
        }
      }).catch(function (err) {
        _this.toast("Delete failure: " + err, "error")
      })
    },

    onToggleModal(toggle) {
      if (toggle) return
      this.onCancelAuthn("DELETE")
      window.clearInterval(this.confirmChallengeInterval)
    },

    onCancelAuthn() {
      this.showConfirm = false
      cancelActivateTerminal(this.syncID, this.cancelSuccess, this.cancelFailure)
    },

    cancelSuccess(resp) { },
    cancelFailure(err) {
      this.toast("错误: " + err, "error")
    },

    onTemporaryAuthn() {
      postAuthorizeTerminal(this.syncID, true, this.authnSuccess, this.authnFailure)
    },

    onConfirmAuthn() {
      postAuthorizeTerminal(this.syncID, false, this.authnSuccess, this.authnFailure)
    },

    authnSuccess(resp) {
      this.showConfirm = false
      this.toast("Success", "success")
    },
    authnFailure(err) {
      this.showConfirm = false
      this.toast("错误: " + err, "error")
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
    document.title = "管理你的终端 -- KnowlGraph"

    // todo 对临时授权的终端进行权限限制

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