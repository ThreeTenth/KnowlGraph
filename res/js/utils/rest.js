// rest.js

function postFinishLogin(terminalID, assertedCredential, success, failure) {
  let authData = new Uint8Array(assertedCredential.response.authenticatorData);
  let clientDataJSON = new Uint8Array(assertedCredential.response.clientDataJSON);
  let rawId = new Uint8Array(assertedCredential.rawId);
  let sig = new Uint8Array(assertedCredential.response.signature);
  let userHandle = new Uint8Array(assertedCredential.response.userHandle);
  axios({
    method: "POST",
    url: queryRestful("/v1/account/t/finishLogin", { id: terminalID }),
    data: JSON.stringify({
      id: assertedCredential.id,
      rawId: bufferEncode(rawId),
      type: assertedCredential.type,
      response: {
        authenticatorData: bufferEncode(authData),
        clientDataJSON: bufferEncode(clientDataJSON),
        signature: bufferEncode(sig),
        userHandle: bufferEncode(userHandle),
      },
    }),
    contentType: "application/json; charset=utf-8",
    dataType: "json",
  }).then(success).catch(failure)
}

function putBeginLogin(terminalID, success, failure) {
  axios({
    method: "PUT",
    url: queryRestful("/v1/account/t/beginLogin", { id: terminalID })
  }).then((resp) => {
    var makeAssertionOptions = resp.data
    makeAssertionOptions.publicKey.challenge = bufferDecode(makeAssertionOptions.publicKey.challenge);
    makeAssertionOptions.publicKey.allowCredentials.forEach(function (listItem) {
      listItem.id = bufferDecode(listItem.id)
    });
    success(resp)
  }).catch(failure)
}

function postFinishRegistration(challenge, state, newCredential, success, failure) {
  let attestationObject = new Uint8Array(newCredential.response.attestationObject);
  let clientDataJSON = new Uint8Array(newCredential.response.clientDataJSON);
  let rawId = new Uint8Array(newCredential.rawId);

  axios({
    method: "POST",
    url: queryRestful("/v1/account/t/finishRegistration", { state: state, challenge: challenge }),
    data: JSON.stringify({
      id: newCredential.id,
      rawId: bufferEncode(rawId),
      type: newCredential.type,
      response: {
        attestationObject: bufferEncode(attestationObject),
        clientDataJSON: bufferEncode(clientDataJSON),
      },
    }),
    contentType: "application/json; charset=utf-8",
    dataType: "json",
  }).then(success).catch(failure)
}

function putBeginRegistration(challenge, state, success, failure) {
  axios({
    method: "PUT",
    url: queryRestful("/v1/account/t/beginRegistration", { state: state, challenge: challenge })
  }).then((resp) => {
    var makeCredentialOptions = resp.data
    makeCredentialOptions.publicKey.challenge = bufferDecode(makeCredentialOptions.publicKey.challenge);
    makeCredentialOptions.publicKey.user.id = bufferDecode(makeCredentialOptions.publicKey.user.id);
    if (makeCredentialOptions.publicKey.excludeCredentials) {
      for (var i = 0; i < makeCredentialOptions.publicKey.excludeCredentials.length; i++) {
        makeCredentialOptions.publicKey.excludeCredentials[i].id = bufferDecode(makeCredentialOptions.publicKey.excludeCredentials[i].id);
      }
    }
    success(resp)
  }).catch(failure)
}

function putMakeCredential(challenge, state, success, failure) {
  axios({
    method: "PUT",
    url: queryRestful("/v1/account/t/makeCredential", { state: state, challenge: challenge })
  }).then(success).catch(failure)
}

function getScanChallenge(challenge, success, failure) {
  axios({
    method: "get",
    url: queryRestful("/v1/account/t/scanChallenge", { challenge: challenge }),
  }).then(success).catch(failure)
}

function postAuthorizeTerminal(challenge, onlyOnce, success, failure) {
  axios({
    method: "post",
    url: queryRestful("/v1/account/t/authorize", { challenge: challenge, onlyOnce: onlyOnce }),
  }).then(success).catch(failure)
}

function cancelActivateTerminal(challenge, success, failure) {
  axios({
    method: "delete",
    url: queryRestful("/v1/account/t/activate", { challenge: challenge }),
  }).then(success).catch(failure)
}

function postActivateTerminal(challenge, success, failure) {
  axios({
    method: "post",
    url: queryRestful("/v1/account/t/activate", { challenge: challenge }),
  }).then(success).catch(failure)
}

function postAuthorizeTerminal(challenge, onlyOnce, success, failure) {
  axios({
    method: "post",
    url: queryRestful("/v1/account/t/authorize", { challenge: challenge, onlyOnce: onlyOnce }),
  }).then(success).catch(failure)
}

/////////////////////////////////////////////////////////////////////////////////

function queryPage(url, data = undefined) {
  url = window.location.origin + url
  return encodeQueryData(url, data)
}

function queryRestful(url, data = undefined) {
  url = restfulDomain + url
  return encodeQueryData(url, data)
}

function queryStatic(url, data = undefined) {
  url = staticDomain + url
  return encodeQueryData(url, data)
}

function encodeQueryData(url, data = undefined) {
  if (undefined == data) return url
  const ret = [];
  for (let d in data)
    ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));

  if (0 == ret.length) return url

  return url + "?" + ret.join('&');
}

function getStatus4Error(err) {
  return err.response ? err.response.status : 9999
}