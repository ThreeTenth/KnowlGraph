// main.js
// Created at 2021-02-02

var storyID;
function onNewStory() {
  axios({
    method: "PUT",
    url: "/api/v1/story",
  }).then(function (resp) {
    console.log(resp.status, resp.data)
    storyID = resp.data
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

var contentTextarea = document.getElementById("content")
var lastContent = contentTextarea.value
var timeoutID;
function onPostChanged() {
  window.clearTimeout(timeoutID);
  timeoutID = window.setTimeout(postStoryContent, 2000);
}

function postStoryContent() {
  if (contentTextarea.value != lastContent) {
    axios({
      method: "PUT",
      url: "/api/v1/story/content?lang=" + getLang(),
      data: {
        content: contentTextarea.value,
        storyID: storyID,
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

function onGetStories() {
  axios({
    method: "GET",
    url: "/api/v1/stories",
  }).then(function (resp) {
    var stories_layout = document.getElementById("stories_layout")
    var stories = resp.data
    for (let index = 0; index < stories.length; index++) {
      const content = stories[index];
      const content_a = document.createElement("a")
      content_a.innerHTML = content.seo.substring(0, 48)
      content_a.href = "#" + content.id
      content_a.onclick = function () {
        editStoryContent(content.story_versions, content.id)
      }
      const content_div = document.createElement("div")
      content_div.appendChild(content_a)
      stories_layout.appendChild(content_div)
    }
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
}

function editStoryContent(_storyID, _contentID) {
  axios({
    method: "GET",
    url: "/api/v1/story?id=" + _storyID + "&content_id=" + _contentID,
  }).then(function (resp) {
    const content = resp.data
    document.getElementById("content_layout").style.display = "block"
    document.querySelector("#langs_select").value = content.edges.Lang.id
    registerLangSelect(document.querySelector("#langs_select"))
    setContentValue(content.content)
    storyID = _storyID
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