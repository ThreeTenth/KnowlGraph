var storyID;
function onNewStory() {
  axios({
    method: "PUT",
    url: "/api/v1/story",
  }).then(function (resp) {
    console.log(resp.status, resp.data)
    storyID = resp.data
    document.getElementById("content_layout").style.display = "block"
    document.getElementById("langs_select").value = getLang()
    contentTextarea.value = ""
  }).catch(function (resp) {
    console.log(resp.status, resp.data)
  })
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
  }
  lastContent = contentTextarea.value
}

function getLang() {
  var language;
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
    stories_layout.value = resp.data
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