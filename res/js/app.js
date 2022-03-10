// app.js
// Created at 03/09/2021

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register("/sw.js", { scope: '/' })
}

let userLang = {
  code: "en",
  name: "English",
}

languages.forEach(element => {
  if (element.code == getUserLang()) {
    userLang = element
    return
  }
});

const i18n = Vue.observable({
  ...defaultLang,
})

var _voteSource = {
  exist: false,
  value: null,
}
const _voteObservable = Vue.observable({
  ..._voteSource,
})

let i18ns = new Map();
i18ns.set(defaultLang.language, defaultLang)

axios.defaults.headers.common['Authorization'] = Cookies.get("access_token")

const plugin = {
  install: function (Vue, options) {
    Vue.prototype.user = {
      lang: userLang,
      logined: logined,
      words: [],
      archives: [],
    }
    Vue.prototype.languages = languages
    Vue.prototype.i18n = i18n
    Vue.prototype.vote = _voteObservable
    Vue.prototype.getUserWords = function () {
      axios({
        method: "GET",
        url: queryRestful("/v1/words"),
      }).then(function (resp) {
        Vue.prototype.user.words = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    }
    Vue.prototype.getArchives = function () {
      // todo 缓存，及归档转换为节点。
      // 含 addArchives 函数
      axios({
        method: "GET",
        url: queryRestful("/v1/archives"),
      }).then(function (resp) {
        Vue.prototype.user.archives = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    }
    Vue.prototype.addArchive = function(archive) {
      var archives = Vue.prototype.user.archives
      archives.forEach(element => {
        if (element.id == archive.id) {
          return
        }
      })
      archives.push(archive)
    }

    Vue.prototype.md2html = function (md) {
      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'strikethrough': 'true',
        'simplifiedAutoLink': 'true',
        'smoothLivePreview': 'true',
        'smartIndentationFix': 'true',
        'openLinksInNewWindow': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      var html = document.createElement('div')
      html.innerHTML = converter.makeHtml(md)

      var ps = html.getElementsByTagName('p')
      for (let index = 0; index < ps.length; index++) {
        const element = ps[index];
        element.classList.add('quote')
      }

      return html.innerHTML
    }

    Vue.prototype.emoji = function (x, i) {
      return '&#x' + (x + i).toString(16)
    }

    Vue.prototype.toast = function (message, status = 'info', duration = 3000) {
      const Toast = Vue.extend(ToastObject)
      const ins = new Toast()
      ins.$mount(document.createElement('div'))
      document.body.appendChild(ins.$el)
      ins.message = message
      ins.status = status
      ins.visible = true
      setTimeout(() => {
        ins.visible = false
      }, duration)
    }

    Vue.prototype.todo = function (msg) {
      Vue.prototype.toast(isString(msg) ? msg : "正在开发中")
    }

    Vue.prototype.back = function () {
      if (1 == history.length) {
        router.push({ path: "/" })
      } else {
        history.back()
      }
    }

    Vue.prototype.canCreateAccount = canUseWenAuthn()
    Vue.prototype.isMobile = detectMob()
    Vue.prototype.isAndroid = detectAndroid()
    Vue.prototype.isIOS = detectIOS()

    const turndownService = new TurndownService({
      codeBlockStyle: "fenced",
      headingStyle: "atx"
    })
    turndownService.use([
      turndownPluginGfm.gfm,
      turndownPluginGfm.tables,
      turndownPluginGfm.strikethrough
    ])
    Vue.prototype.turndownService = turndownService

    Vue.prototype.onGitHubOAuth = function () {
      const github_client_id = getMeta("github_client_id")
      const state = Math.random().toString(36).slice(2)
      const githubOAuthAPI = "https://github.com/login/oauth/authorize?client_id=" + github_client_id + "&state=" + state
      Cookies.set('github_oauth_state', state)
      window.open(githubOAuthAPI, "_self")
    }
  }
}

// 注册一个全局自定义指令 `v-focus`
Vue.directive('focus', (el, binding) => {
  if (binding.value) {
    el.focus()
  }
})

const About = { template: fgm_about }
const My = { template: fgm_my }
const GetAPP = { template: fgm_get_app }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', name: "index", component: Index },
    { path: '/p/:id/:code?', name: "article", component: Article, props: true },
    { path: '/drafts', name: 'drafts', component: Drafts },
    { path: '/archive', component: Archive, children: archive_router },
    { path: '/about', name: 'about', component: About },
    { path: '/my', name: 'my', component: My },
    { path: '/settings/personalize', name: 'personalize', component: Personalize },
    { path: '/terminals', name: 'terminals', component: Terminals },
    { path: '/new/article', name: 'newArticle', component: NewArticle },
    { path: '/d/:id/edit', name: 'editDraft', component: EditDraft, props: true },
    { path: '/login', name: 'login', component: NewAccount },
    { path: '/g/:id', redirect: { name: 'getapp' } },
    { path: '/getapp', name: 'getapp', component: GetAPP },
    { path: '/dashboard', name: 'dashboard', component: Dashboard },
    { path: '/vote', name: 'vote', component: Vote },
    { path: '/word/:id/:name?', name: 'word', component: Word, props: true },
    { path: '/words', name: 'words', component: Words },
    { path: '/node/:id/:name?', name: 'node', component: NodeArticles, props: true },
    { path: '/nodes/:id?/:name?', name: 'nodes', component: Nodes, props: true },
  ]
})

router.beforeEach((to, from, next) => {
  // ...
  var page = to.fullPath
  var title = to.name
  if (to.name === "article") {
    let last = to.fullPath.lastIndexOf("/")
    page = to.fullPath.substring(0, last)
    title = to.fullPath.substr(last + 1)
    title = decodeURIComponent(title)
  } else if (to.name === "editDraft") {
    page = "/d"
    title = "edit draft"
  } else if (to.name === "newArticle") {
    page = "/new/article"
    title = "new article"
  }
  next()
  postAnalyticsPageView(page, title, from.fullPath)
})

var app = new Vue({
  data: {},
  router,
  // template: logined ? app_home : app_index,

  created() {
    if (logined) {
      axios({
        method: "GET",
        url: queryRestful("/v1/vote"),
      }).then((resp) => {
        if (200 == resp.status) {
          Object.assign(_voteObservable, { exist: true, value: resp.data })
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    }
  },

  methods: {
    onSignout: function () {
      const githubOAuthAPI = "/signout"
      window.open(githubOAuthAPI, "_self")
    },
  }
})

function runApp() {
  Vue.use(plugin)
  app.$mount('#app')

  if (logined) {
    app.getUserWords()
    app.getArchives()
  }
}

var langCode = getUserLang()
if (langCode == defaultLang.__code) {
  runApp()
} else {
  var url = queryStatic("/strings/strings-" + langCode + ".json")

  axios({
    method: "GET",
    url: url,
  }).then(function (resp) {
    i18ns.set(langCode, resp.data)
    setI18nStrings(langCode, resp.data)
    Object.assign(i18n, resp.data)

    runApp()
  }).catch(function (resp) {
    console.log(resp)
  })
}
