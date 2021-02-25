// main.js
// Created at 2021-02-02

getTags()

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

var postChangedTimeoutID;
var cancelDeleteTagTimeoutID;

var lastContent;

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
    seen: true,
    lang: getLang(),
  },
  methods: {
    onSelectLang: function (event) {
      lang_layout.lang = event.target.value
    },
  }
})

var content_layout = new Vue({
  el: '#content_layout',
  data: {
    seen: false,
    body: "",
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
    auto: false,
    tags: [],
    tag: "",
    items: [],
    all: [],
    state: TagStateMayDel,
  },
  methods: {
    onTagChanged: function () {
      tagChanged()
    },
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

function postChanged() {
  window.clearTimeout(postChangedTimeoutID)
  postChangedTimeoutID = window.setTimeout(postArticleContent, 800)
}

function postBlur() {
  window.clearTimeout(postChangedTimeoutID)
  postArticleContent()
}

function tagChanged() {
  if ("" == tag_layout.tag) {
    if (TagStateSafe == tag_layout.state) {
      tag_layout.state = TagStateIdle
    }
    tag_layout.auto = false
  } else {
    tag_layout.state = TagStateSafe
    if (2 <= tag_layout.tag.length) {
      const items = []
      const regex = new RegExp('^' + tag_layout.tag)
      tag_layout.all.forEach(item => {
        if (regex.test(item)) {
          items.push(item)
        }
      });
      tag_layout.items = items
      tag_layout.auto = true
    } else {
      tag_layout.auto = false
    }
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

function setContent(content = undefined) {
  if (content == undefined) {
    content_layout.body = ""
    tag_layout.tags = []
  } else {
    content_layout.body = content.body

    var tags = []
    content.edges.Tags.forEach(tag => {
      tags.push(tag.name)
    });
    tag_layout.tags = tags
  }
  content_layout.seen = true
  tag_layout.seen = true
  tag_autocomplete.seen = true

  lastContent = getCurrentContent()
}

function getCurrentContent() {
  return content_layout.body + tag_layout.tags.join(",")
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

function encodeQueryData(url, data = undefined) {
  if (undefined == data) return url
  const ret = [];
  for (let d in data)
    ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));

  if (0 == ret.length) return url

  return url + "?" + ret.join('&');
}
