// main.js
// Created at 2021-02-02

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

var postChangedTimeoutID;
var cancelDeleteTagTimeoutID;

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
    lang: getLang(),
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
  },
  methods: {
    onChanged: function () {
      postChanged()
    },
    onBlur: function () {
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
    state: TagStateMayDel,
  },
  methods: {
    onRemoveTag: function (tag) {
      removeTag(tag)
    },
    onAdd: function () {
      addTag()
    },
    onDel: function () {
      deleteTag()
    },
  },
})

var lastContent;
function postChanged() {
  window.clearTimeout(postChangedTimeoutID)
  postChangedTimeoutID = window.setTimeout(postArticleContent, 800)
}

function postBlur() {
  window.clearTimeout(postChangedTimeoutID)
  postArticleContent()
}

function postArticleContent() {
  var currentContent = getCurrentContent()
  if (currentContent != lastContent) {
    axios({
      method: "PUT",
      url: "/api/v1/article/content?lang=" + getLang(),
      data: {
        content: content_layout.content,
        articleID: articleID,
        tags: tag_layout.tags,
      },
    }).then(function (resp) {
      console.log(resp.status, resp.data)
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
    lastContent = currentContent
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

function tagChanged() {
  if ("" == tag_layout.tag) {
    if (TagStateSafe == tag_layout.state) {
      tag_layout.state = TagStateIdle
    }
  } else {
    tag_layout.state = TagStateSafe
  }
}

function removeTag(tag) {
  const index = tag_layout.tags.indexOf(tag)
  if (index > -1) {
    tag_layout.tags.splice(index, 1)
  }
  tag_layout.state = TagStateMayDel

  postBlur()
}

function addTag() {
  // var tag = tag_layout.tag.trim()         // ",vue,,,,,   ,web, golang      ,,  ,,rest,"
  // tag = tag.replace(/\, +,/g, '')         // ",vue,,,,web, golang      ,,rest,"
  // tag = tag.replace(/\,+/g, ',')          // ",vue,web, golang      ,rest,"
  // tag = tag.replace(/^\,+|\,+$/g, '')     // "vue,web, golang      ,rest"
  // tag = tag.replace(/ *\, */g, ',')       // "vue,web,golang,rest"
  //                                         // "gin    ", "gin  ,", "gin,", "gin,， ，  ,"
  // tag = tag.replace(/ *\,*$/g, '')        // "gin"

  var tag = tag_layout.tag.replace(/[\,\， ]*$/g, '')
  if ("" == tag) {
    tag_layout.tag = ""
    return
  }

  tag_layout.tags.push(tag)
  tag_layout.tag = ""
  tag_layout.state = TagStateMayDel

  postBlur()
}

function deleteTag() {
  switch (tag_layout.state) {
    case TagStateIdle:
      tag_layout.state = TagStateMayDel
      break;
    case TagStateMayDel:
      tag_layout.state = TagStatePreDel
      cancelDeleteTagTimeoutID = window.setTimeout(cancelDeleteTag, 2000)
      break;
    case TagStatePreDel:
      const removeIndex = tag_layout.tags.length
      tag_layout.tags.splice(removeIndex - 1, 1)
      tag_layout.state = TagStateMayDel
      window.clearTimeout(cancelDeleteTagTimeoutID)

      postBlur()
      break;
  }
}

function cancelDeleteTag() {
  tag_layout.state = TagStateMayDel
}

function setLang(lang) {
  Cookies.set("user-lang", lang)
  content_layout.content = ""
  lastContent = getCurrentContent()
}

function setContent(content = undefined) {
  if (content == undefined) {
    content_layout.content = ""
    lang_layout.lang = getLang()
    tag_layout.tags = []
  } else {
    content_layout.content = content.content
    lang_layout.lang = content.edges.Lang.id

    var tags = []
    content.edges.Tags.forEach(tag => {
      tags.push(tag.name)
    });
    tag_layout.tags = tags
  }
  content_layout.seen = true
  lang_layout.seen = true
  tag_layout.seen = true

  lastContent = getCurrentContent()
}

function getCurrentContent() {
  return content_layout.content + tag_layout.tags.join(",")
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