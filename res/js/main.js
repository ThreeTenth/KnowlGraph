// main.js
// Created at 2021-02-02

var articleID;
function onNewArticle() {
  axios({
    method: "PUT",
    url: "/api/v1/article",
  }).then(function (resp) {
    console.log(resp.status, resp.data)
    articleID = resp.data
    document.getElementById("content_layout").style.display = "block"
    document.querySelector("#langs_select").value = getLang()
    registerLangSelect(document.querySelector("#langs_select"))
    setContentValue("")
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function registerLangSelect(selectElement) {
  selectElement.addEventListener('change', (event) => {
    Cookies.set("user-lang", event.target.value)
    setContentValue("")
  });
}

var tagInput = document.getElementById("tag_input")
var contentTextarea = document.getElementById("content")
var lastContent = contentTextarea.value
var timeoutID;
function onPostChanged() {
  window.clearTimeout(timeoutID);
  timeoutID = window.setTimeout(postArticleContent, 2000);

  var tagString = tagInput.value
  console.log(tagString)
}

function postArticleContent() {
  if (contentTextarea.value != lastContent) {
    axios({
      method: "PUT",
      url: "/api/v1/article/content?lang=" + getLang(),
      data: {
        content: contentTextarea.value,
        articleID: articleID,
      },
    }).then(function (resp) {
      console.log(resp.status, resp.data)
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
    lastContent = contentTextarea.value
  }
}

function getLang() {
  var language = Cookies.get("user-lang");
  if (language != undefined)
    return language

  if (navigator.languages != undefined)
    language = navigator.languages[0];
  else
    language = navigator.language || window.navigator.userLanguage;

  var localLang = language.split("-")
  if (2 == localLang.length)
    language = localLang[0]

  return language
}

function getMeta(metaName) {
  const metas = document.getElementsByTagName('meta');

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('name') === metaName) {
      return metas[i].getAttribute('content');
    }
  }

  return '';
}

function onGetArticles() {
  axios({
    method: "GET",
    url: "/api/v1/articles",
  }).then(function (resp) {
    var articles_layout = document.getElementById("articles_layout")
    var articles = resp.data

    articles_layout.innerHTML = ""
    for (let index = 0; index < articles.length; index++) {
      const content = articles[index];
      const content_a = document.createElement("a")
      content_a.innerHTML = content.seo.substring(0, 48)
      content_a.href = "#" + content.id
      content_a.onclick = function () {
        editArticleContent(content.article_versions, content.id)
      }
      const content_div = document.createElement("div")
      content_div.appendChild(content_a)
      articles_layout.appendChild(content_div)
    }
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function editArticleContent(_articleID, _contentID) {
  axios({
    method: "GET",
    url: "/api/v1/article?id=" + _articleID + "&content_id=" + _contentID,
  }).then(function (resp) {
    const content = resp.data
    document.getElementById("content_layout").style.display = "block"
    document.querySelector("#langs_select").value = content.edges.Lang.id
    registerLangSelect(document.querySelector("#langs_select"))
    setContentValue(content.content)
    articleID = _articleID
  }).catch(function (resp) {

  })
}

function setContentValue(value) {
  contentTextarea.value = value
  lastContent = value
}

function onSignout() {
  const githubOAuthAPI = "/signout"
  window.open(githubOAuthAPI, "_self")
}

function onGitHubOAuth() {
  const github_client_id = getMeta("github_client_id")
  const state = Math.random().toString(36).slice(2)
  const githubOAuthAPI = "https://github.com/login/oauth/authorize?client_id=" + github_client_id + "&state=" + state
  Cookies.set('github_oauth_state', state)
  window.open(githubOAuthAPI, "_self")
}