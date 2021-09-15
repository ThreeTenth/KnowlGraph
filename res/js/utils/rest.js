// rest.js

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