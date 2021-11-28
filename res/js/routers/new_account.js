// new_account.js

const NewAccount = {

  data: function () {
    return {
      requestState: 0,
      challenge: "",
      challengeState: 0,
      isQRCode: true,
    }
  },

  methods: {

    switchAuthMethod() {
      var is = this.isQRCode
      this.isQRCode = !is
    },

    onChallenge(challenge, requestState) {
      this.challenge = challenge
      this.requestState = requestState
    },

    onResult(qrcodeResult) {
      var state = qrcodeResult.state
      if (this.challengeState == state) {
        return
      }

      if (state == 4) {
        this.makeCredential(this.requestState, this.challenge)
      }

      this.challengeState = state
    },

    makeCredential(state, challenge) {
      // axios({
      //   method: "PUT",
      //   url: queryRestful("/v1/account/t/makeCredential", { state: state, challenge: challenge })
      // }).then(function (resp) {
      //   authSuccess(resp.data)
      // }).catch(function (err) { })
      putBeginRegistration(challenge, state, (resp) => {
        this.requestWebAuthn(state, challenge, resp.data)
      }, this.beginRegistrationFailure)
    },

    beginRegistrationFailure(err) {
      this.toast(err, "error")
    },
    requestWebAuthn(state, challenge, makeCredentialOptions) {
      makeCredentialOptions.publicKey.challenge = bufferDecode(makeCredentialOptions.publicKey.challenge);
      makeCredentialOptions.publicKey.user.id = bufferDecode(makeCredentialOptions.publicKey.user.id);
      if (makeCredentialOptions.publicKey.excludeCredentials) {
        for (var i = 0; i < makeCredentialOptions.publicKey.excludeCredentials.length; i++) {
          makeCredentialOptions.publicKey.excludeCredentials[i].id = bufferDecode(makeCredentialOptions.publicKey.excludeCredentials[i].id);
        }
      }
      // console.log("Credential Creation Options");
      // console.log(makeCredentialOptions);
      navigator.credentials.create(makeCredentialOptions)
        .then((newCredential) => {
          console.log("PublicKeyCredential Created");
          console.log(newCredential);
          postFinishRegistration(challenge, state, newCredential, this.requestWebAuthnSuccess, this.requestWebAuthnFailure)
        }).catch((err) => {
          console.info(err);
          this.toast(err, "error")
        });
    },

    requestWebAuthnSuccess(resp) {
      this.toast("Success", "success")
      console.log(resp.data);
    },
    requestWebAuthnFailure(err) {
      this.toast("Error " + err, "error")
      console.info(err);
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
  },

  template: fgm_new_account,
}