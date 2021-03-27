// main.js
// Created at 2021-02-02

const app = new Vue({ router }).$mount('#app')

var app = new Vue({
  el: "#app",
  data: {
    menu: {
      user: {
        lang: getUserLang(),
      }
    },
    home: {},
    my: {},
    article: {},
    draft: {
      seen: false,
      body: "",
    },
    publish: {
      seen: false,
      lang: getUserLang(),
      words: {
        auto: false,
        values: [],
        new: "",
        items: [],
        all: [],
        state: WordStateMayDel,
      }
    },
    nodes: {},
    toast: [],
  },
  methods: {
    onSelectContentLang: function (event) {
      lang_layout.lang = event.target.value
    },
    onChanged: function () {
      postChanged()
      console.log(i)
    },
    onBlur: function () {
      postBlur()
    },
    onWordChanged: function () {
      wordChanged()
    },
    onRemoveWord: function (word) {
      removeWord(word)
    },
    onAdd: function () {
      addWord()
    },
    onDel: function () {
      deleteWord()
    },
  },
})

var postChangedTimeoutID;
var cancelDeleteWordTimeoutID;

var lastContent;

function postChanged() {
  window.clearTimeout(postChangedTimeoutID)
  postChangedTimeoutID = window.setTimeout(postArticleContent, 800)
}

function postBlur() {
  window.clearTimeout(postChangedTimeoutID)
  postArticleContent()
}

function setContent(content = undefined) {
  if (content == undefined) {
    content_layout.body = ""
  } else {
    content_layout.body = content.body
  }
  content_layout.seen = true

  lastContent = content_layout.body
}

function postArticleContent() {
  if (lastContent == app.draft.body) return
  axios({
    method: "PUT",
    url: encodeQueryData("/api/v1/article/content"),
    data: {
      body: app.draft.body,
      draft_id: draftID,
    },
  }).then(function (resp) {
    console.log(resp.status, resp.data)
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })

  lastContent = content_layout.body
}

function onSelectUserLang() {
  var select = document.getElementById("userLangSelect")
  Cookies.set("user-lang", select.value)
  location.reload()
}

function onNewArticle() {
  axios({
    method: "PUT",
    url: encodeQueryData("/api/v1/article", { status: "private" }),
  }).then(function (resp) {
    console.log(resp.status, resp.data)
    articleID = resp.data.ArticleID
    draftID = resp.data.DraftID
    setContent()
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function onPublishArticle() {
  axios({
    method: "PUT",
    url: encodeQueryData("/api/v1/article/content"),
    data: {
      name: content_layout.body,
      comment: draftID,
    },
  }).then(function (resp) {
    console.log(resp.status, resp.data)
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function onGetArticles() {
  axios({
    method: "GET",
    url: encodeQueryData("/api/v1/user/articles", { status: "self" }),
  }).then(function (resp) {
    articles_layout.seen = true
    articles_layout.articles = resp.data
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function onGetDrafts() {
  axios({
    method: "GET",
    url: encodeQueryData("/api/v1/drafts"),
  }).then(function (resp) {
    articles_layout.seen = true
    articles_layout.articles = resp.data
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
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