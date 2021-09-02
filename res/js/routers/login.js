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
    removeWebauthnSession() {
      Cookies.remove('webauthn-session')
    },
    createAccount() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/challenge"),
      }).then(function (resp) {
        _this.createCredential(resp.data)
      }).catch(function (err) {
        _this.toast("创建请求失败", "error")
        _this.removeWebauthnSession()
      })
    },

    createCredential(credentialCreationOptions) {
      // console.log(credentialCreationOptions)
      credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
      credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
      if (credentialCreationOptions.publicKey.excludeCredentials) {
        for (var i = 0; i < credentialCreationOptions.publicKey.excludeCredentials.length; i++) {
          credentialCreationOptions.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.publicKey.excludeCredentials[i].id);
        }
      }

      credentialCreationOptions.publicKey.user.displayName = ''

      var _this = this
      navigator.credentials.create({ publicKey: credentialCreationOptions.publicKey })
        .then(function (newCredential) {
          // console.log(newCredential)
          let id = newCredential.id;
          let type = newCredential.type;
          let rawId = newCredential.rawId;
          let clientDataJSON = newCredential.response.clientDataJSON;
          let attestationObject = newCredential.response.attestationObject;

          attestationObject = new Uint8Array(newCredential.response.attestationObject);
          clientDataJSON = new Uint8Array(newCredential.response.clientDataJSON);
          rawId = new Uint8Array(newCredential.rawId);

          _this.putFinishCreate(id, type, bufferEncode(rawId),
            bufferEncode(clientDataJSON), bufferEncode(attestationObject))
        }).catch(function (err) {
          console.error(err)
          _this.toast("用户操作超时或拒绝授权", "error")
          _this.removeWebauthnSession()
        });
    },

    putFinishCreate(credentialId, credentialType, rawId, clientDataJSON, attestationObject) {
      var _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/account/create"),
        data: {
          id: credentialId,
          rawId: rawId,
          type: credentialType,
          response: {
            attestationObject: attestationObject,
            clientDataJSON: clientDataJSON,
          },
        },
      }).then(function (resp) {
        _this.toast("授权成功", "success")
      }).catch(function (err) {
        _this.toast("授权失败", "error")
        _this.removeWebauthnSession()
      })
    },

    authnAccount() {
      let _this = this;
      axios({
        method: "GET",
        url: queryRestful("/v1/account/request"),
      }).then(function (resp) {
        _this.getCredentials(resp.data)
      }).catch(function (err) {
        _this.toast("授权失败", "error")
        _this.removeWebauthnSession()
      })
    },

    async getCredentials(credentialRequestOptions) {
      // console.log(credentialRequestOptions)
      credentialRequestOptions.publicKey.challenge = bufferDecode(credentialRequestOptions.publicKey.challenge);
      credentialRequestOptions.publicKey.allowCredentials.forEach(function (listItem) {
        listItem.id = bufferDecode(listItem.id)
      });

      var _this = this
      navigator.credentials.get({ publicKey: credentialRequestOptions.publicKey })
        .then(function (assertion) {
          _this.postCredentials(assertion)
        })
        .catch(function (err) {
          console.error(err)
          _this.toast("用户操作超时或拒绝授权", "error")
          _this.removeWebauthnSession()
        });
    },

    postCredentials(assertion) {
      let authData = assertion.response.authenticatorData;
      let clientDataJSON = assertion.response.clientDataJSON;
      let rawId = assertion.rawId;
      let sig = assertion.response.signature;
      let userHandle = assertion.response.userHandle;

      var _this = this
      axios({
        method: "POST",
        url: queryRestful("/v1/account/anthn"),
        data: {
          id: assertion.id,
          rawId: bufferEncode(new Uint8Array(rawId)),
          type: assertion.type,
          response: {
            authenticatorData: bufferEncode(new Uint8Array(authData)),
            clientDataJSON: bufferEncode(new Uint8Array(clientDataJSON)),
            signature: bufferEncode(new Uint8Array(sig)),
            userHandle: bufferEncode(new Uint8Array(userHandle)),
          },
        },
      }).then(function (resp) {
        _this.toast("授权成功", "success")
      }).catch(function (err) {
        _this.toast("授权失败", "error")
        _this.removeWebauthnSession()
      })
    },
  },

  created() {
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