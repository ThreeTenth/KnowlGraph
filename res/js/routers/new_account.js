// new_account.js

const NewAccount = {

  data: function () {
    return {
      requestState: 0,
      challenge: "",
      challengeState: 0,
      isQRCode: true,
      isExpiredTerminal: false,
      isSwitchTemporaryAuth: false,
    }
  },

  computed: {
    terminalName: function () {
      return Cookies.get("terminal_name")
    }
  },

  methods: {

    onCancelTemporaryAuth() {
      this.isSwitchTemporaryAuth = false
    },

    onSwitchTemporaryAuth() {
      postFinishRegOnlyOnce(this.challenge, this.requestState,
        this.finishRegOnlyOnceSuccess, this.finishRegOnlyOnceFailure)
    },

    backOrWebAuthn() {
      if (isExpiredTerminal()) {
        this.isExpiredTerminal = true
      } else {
        this.back()
      }
    },

    startWebAuthn() {
      var tid = Cookies.get("terminal_id")
      putBeginLogin(tid, (resp) => {
        this.beginLoginSuccess(tid, resp.data)
      }, this.beginLoginFailure)
    },

    beginLoginFailure(err) {
      this.toast(err, "error")
      console.info("BeginLoginFailure", err)
    },
    beginLoginSuccess(terminalID, makeAssertionOptions) {
      console.log("Assertion Options:");
      // console.log(makeAssertionOptions);
      navigator.credentials.get(makeAssertionOptions)
        .then((credential) => {
          // console.log(credential);
          postFinishLogin(terminalID, credential, this.finishLoginSuccess, this.finishLoginFailure)
        }).catch((err) => {
          this.toast(err, "error")
          console.info("Error", err)
        });
    },

    finishLoginFailure(err) {
      this.toast(err, "error")
      console.info("FinishLoginFailure", err)
    },
    finishLoginSuccess(resp) {
      // console.log(resp.data)
      authSuccess(resp.data)
    },

    switchOtherAccount() {
      this.isExpiredTerminal = false
    },

    switchAuthMethod() {
      var is = this.isQRCode
      this.isQRCode = !is
    },

    onChallenge(challenge, requestState) {
      this.challenge = challenge
      this.requestState = requestState
    },

    onResult(qrcodeResult) {
      // var name = qrcodeResult.name
      var state = qrcodeResult.state
      var onlyOnce = qrcodeResult.onlyOnce
      if (this.challengeState == state) {
        return
      }

      if (state == 4) {
        if (onlyOnce) {
          // 临时授权流程
          postFinishRegOnlyOnce(this.challenge, this.requestState,
            this.finishRegOnlyOnceSuccess, this.finishRegOnlyOnceFailure)
        } else if (canUseWenAuthn()) {
          this.makeCredential(state, this.requestState, this.challenge)
        } else {
          // 长期授权，但终端没有 WebAuthn 认证功能时
          this.isSwitchTemporaryAuth = true
        }
      }

      this.challengeState = state
    },

    makeCredential(authState, state, challenge) {
      putBeginRegistration(challenge, state, (resp) => {
        this.requestWebAuthn(authState, state, challenge, resp.data)
      }, this.beginRegistrationFailure)
    },

    beginRegistrationFailure(err) {
      this.toast(err, "error")
    },
    requestWebAuthn(authState, state, challenge, makeCredentialOptions) {
      console.log("Credential Creation Options");
      // console.log(makeCredentialOptions);
      navigator.credentials.create(makeCredentialOptions)
        .then((newCredential) => {
          console.log("PublicKeyCredential Created");
          // console.log(newCredential);
          postFinishRegistration(challenge, state, newCredential, this.requestWebAuthnSuccess, this.requestWebAuthnFailure)
        }).catch((err) => {
          console.info(err);
          this.requestWebAuthnError(authState, state, challenge)
        });
    },
    requestWebAuthnError(authState, state, challenge) {
      if (4 != authState) {
        this.toast("认证失败", "error")
        return
      }

      // 如果是用户添加终端，但终端不支持(或用户无法认证，此条待确认) WebAuthn，则为用户提供切换到临时授权模式选项。
      // todo 或用户无法认证，此条待确认
      this.isSwitchTemporaryAuth = true
      // postFinishRegOnlyOnce(challenge, state, this.finishRegOnlyOnceSuccess, this.finishRegOnlyOnceFailure)
    },
    finishRegOnlyOnceFailure(err) {
      this.toast("Error " + err, "error")
      console.info(err);
    },
    finishRegOnlyOnceSuccess(resp) {
      this.toast("Success", "success")
      authSuccess(resp.data)
    },

    requestWebAuthnSuccess(resp) {
      this.toast("Success", "success")
      authSuccess(resp.data)
    },
    requestWebAuthnFailure(err) {
      this.toast("Error " + err, "error")
      console.info(err);
    },

    createAccount() {
      let state = Math.random().toString(36).slice(2)
      axios({
        method: "GET",
        url: queryRestful("/v1/account/t/challenge", { state: state }),
      }).then((resp) => {
        this.makeCredential(1, state, resp.data)
      }).catch((err) => {
        this.toast(err, "error")
      })
    },
  },

  created() {
    document.title = "同步账号 -- KnowlGraph"

    this.isExpiredTerminal = isExpiredTerminal()
  },

  beforeDestroy() {
  },

  template: fgm_new_account,
}