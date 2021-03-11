// utils.js

function getUserLang() {
  var language = Cookies.get("user-lang")
  if (language != undefined)
    return language

  if (navigator.languages != undefined)
    language = navigator.languages[0]
  else
    language = navigator.language || window.navigator.userLanguage

  var localLang = language.split("-")
  if (2 == localLang.length)
    language = localLang[0]

  return language
}

function setUserLang(lang) {
  Cookies.set("user-lang", lang)
}

function getMeta(metaName) {
  const metas = document.getElementsByTagName('meta')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('name') === metaName) {
      return metas[i].getAttribute('content')
    }
  }

  return ''
}

function setI18nStrings(lang, strings) {
  localStorage.setItem('strings-' + lang, JSON.stringify(strings));
}

function getI18nStrings(lang) {
  let strings = localStorage.getItem('strings-' + lang);
  return JSON.parse(strings)
}

function removeI18nStrings(lang) {
  localStorage.removeItem('strings-' + lang);
}