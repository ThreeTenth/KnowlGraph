// app.js
// Created at 03/09/2021

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

const plugin = {
  install: function (Vue, options) {
    Vue.prototype.user = {
      lang: userLang,
      logined: logined,
      words: [],
      archives: [],
    }
    Vue.prototype.i18n = i18n
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
      axios({
        method: "GET",
        url: queryRestful("/v1/archives"),
      }).then(function (resp) {
        Vue.prototype.user.archives = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    }

    Vue.prototype.md2html = function (md) {
      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      return converter.makeHtml(md);
    }
  }
}

Vue.use(plugin)

const About = { template: fgm_about }
const My = { template: fgm_my }
const Login = { template: fgm_login }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', component: Index },
    { path: '/p/:id/:code?', name: "article", component: Article, props: true },
    { path: '/drafts', name: 'drafts', component: Drafts },
    { path: '/archive', component: Archive, children: archive_router },
    { path: '/about', name: 'about', component: About },
    { path: '/my', name: 'my', component: My },
    { path: '/new/article', name: 'newArticle', component: NewArticle },
    { path: '/d/:id/edit', name: 'editDraft', component: EditDraft, props: true },
    { path: '/d/:id/history', name: 'draftHistories', component: DraftHistories, props: true },
    { path: '/d/:id/history/:hid', name: 'draftHistory', component: DraftHistory, props: true },
    { path: '/d/:id/publish', name: 'publishArticle', component: PublishArticle, props: true },
    { path: '/login', name: 'login', component: Login },
  ]
})

var app = new Vue({
  data: {
    vote: {
      has: false,
      ras: null,
    },
    languages: languages,
    profilePicture: getLink("icon"),
  },
  router,
  template: logined ? app_home : app_index,

  created() {
    var _this = this
    if (logined) {
      axios({
        method: "GET",
        url: queryRestful("/v1/vote"),
      }).then(function (resp) {
        if (200 == resp.status) {
          _this.vote.ras = resp.data
          _this.vote.has = true
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    }

    this.getUserWords()
    this.getArchives()
  },

  methods: {
    onAllow() {
      if (!this.ras) return

      this.__postVote("allowed")
    },
    onRejecte() {
      if (!this.ras) return

      this.__postVote("rejected")
    },
    onGitHubOAuth: function () {
      const github_client_id = getMeta("github_client_id")
      const state = Math.random().toString(36).slice(2)
      const githubOAuthAPI = "https://github.com/login/oauth/authorize?client_id=" + github_client_id + "&state=" + state
      Cookies.set('github_oauth_state', state)
      window.open(githubOAuthAPI, "_self")
    },
    onSignout: function () {
      const githubOAuthAPI = "/signout"
      window.open(githubOAuthAPI, "_self")
    },
    onSelectUserLang: function (lang) {
      // const lang = event.target.value
      const old = this.user.lang
      const langCode = lang.code

      this.user.lang = lang
      setUserLang(langCode)

      if (langCode == defaultLang.__code) {
        Object.assign(i18n, defaultLang)
        return
      } else if (i18ns.has(langCode)) {
        Object.assign(i18n, i18ns.get(langCode))
        return
      } else {
        let strings = getI18nStrings(langCode)
        if (null != strings && strings.__version == defaultLang.__version) {
          i18ns.set(langCode, strings)
          Object.assign(i18n, strings)
          return
        }
        removeI18nStrings(langCode)
      }

      let _this = this

      axios({
        method: "GET",
        url: queryStatic("/static/strings/strings-" + langCode + ".json"),
      }).then(function (resp) {
        i18ns.set(langCode, resp.data)
        setI18nStrings(langCode, resp.data)
        Object.assign(i18n, resp.data)
      }).catch(function (resp) {
        _this.user.lang = old
        setUserLang(old.code)
      })
    },
    __postVote(status) {
      var _this = this
      axios({
        method: "POST",
        url: queryRestful("/v1/vote"),
        data: {
          id: this.ras.id,
          status: status,
        },
      }).then(function (resp) {
        _this.ras = null
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  }
}).$mount('#application--wrap')