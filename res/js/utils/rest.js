// rest.js

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