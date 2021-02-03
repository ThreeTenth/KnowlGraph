// main.js
// Created at 2021-02-02

var articles_layout = new Vue({
  el: '#articles_layout',
  data: {
    seen: false,
    articles: [],
  },
  methods: {
    onEditArticleContent: function (_articleID, _contentID) {
      editArticleContent(_articleID, _contentID)
    },
  }
})

var lang_layout = new Vue({
  el: '#lang_layout',
  data: {
    seen: false,
    default_lang: getLang(),
  },
  methods: {
    onSelectLang: function (event) {
      setLang(event.target.value)
    },
  }
})

var content_layout = new Vue({
  el: '#content_layout',
  data: {
    seen: false,
    content: "",
    default_lang: getLang(),
  },
  methods: {
    onSelectLang: function (event) {
      setLang(event.target.value)
    },
    onPostChanged: function () {
      postChanged()
    },
    onPostBlur: function () {
      postBlur()
    },
  }
})

var tag_layout = new Vue({
  el: '#tag_layout',
  data: {
    seen: false,
    tags: [],
    tag: "",
  },
  methods: {
    onPostChanged: function () {
      postChanged()
    },
    onPostBlur: function () {
      postBlur()
    },
  },
})

var lastContent;
var timeoutID;
function postChanged() {
  window.clearTimeout(timeoutID)
  timeoutID = window.setTimeout(postArticleContent, 2000)
}

function postBlur() {
  window.clearTimeout(timeoutID)
  postArticleContent()
}

function postArticleContent() {
  if (content_layout.content != lastContent) {
    // const tags = content_layout.tag.split(",")
    axios({
      method: "PUT",
      url: "/api/v1/article/content?lang=" + getLang(),
      data: {
        content: content_layout.content,
        articleID: articleID,
        // tags: tags,
      },
    }).then(function (resp) {
      console.log(resp.status, resp.data)
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
    setLastContent(content_layout.content)
  }
}

function editArticleContent(_articleID, _contentID) {
  axios({
    method: "GET",
    url: "/api/v1/article?id=" + _articleID + "&content_id=" + _contentID,
  }).then(function (resp) {
    const content = resp.data
    articleID = _articleID
    setContent(content)
  }).catch(function (resp) {

  })
}

function setLang(lang) {
  Cookies.set("user-lang", lang)
  setContent()
}

function setContent(content = undefined) {
  if (content == undefined) {
    content_layout.content = ""
    lang_layout.default_lang = getLang()
  } else {
    content_layout.content = content.content
    lang_layout.default_lang = content.edges.Lang.id
  }
  content_layout.seen = true
  lang_layout.seen = true
  tag_layout.seen = true

  setLastContent(content_layout.content)
}

function setLastContent(content) {
  lastContent = content
}

function getContentText() {
  return content_layout.content
}

var articleID;
function onNewArticle() {
  axios({
    method: "PUT",
    url: "/api/v1/article",
  }).then(function (resp) {
    console.log(resp.status, resp.data)
    articleID = resp.data
    setContent()
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function onGetArticles() {
  axios({
    method: "GET",
    url: "/api/v1/articles",
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

function getLang() {
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

function getMeta(metaName) {
  const metas = document.getElementsByTagName('meta')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('name') === metaName) {
      return metas[i].getAttribute('content')
    }
  }

  return ''
}