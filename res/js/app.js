// app.js
// Created at 03/09/2021

const plugin = {
  install: function (Vue, options) {
    Vue.prototype.i18n = i18n
  }
}

Vue.use(plugin)

const Index = { template: logined ? fgm_home : fgm_index }
const Home = { template: fgm_home }
const Drafts = { template: fgm_user_drafts }
const About = { template: fgm_about }
const My = { template: fgm_my }
const Login = { template: fgm_login }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', component: Index },
    { path: '/drafts', name: 'drafts', component: Drafts },
    { path: '/about', name: 'about', component: About },
    { path: '/my', name: 'my', component: My },
    { path: '/new/article', name: 'newArticle', component: NewArticle },
    { path: '/p/:id/edit', name: 'editDraft', component: EditDraft, props: true },
    { path: '/p/:id/histories', name: 'draftHistories', component: DraftHistories, props: true },
    { path: '/p/:id/history/:hid', name: 'draftHistory', component: DraftHistory, props: true },
    { path: '/login', name: 'login', component: Login },
  ]
})

var app = new Vue({
  data: {
    user: {
      lang: getUserLang(),
    },
    languages: languages,
    logined: logined,
  },
  router,
  template: logined ? app_home : app_index,
  methods: {
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
    onSelectUserLang: function (event) {
      const lang = event.target.value
      const old = this.user.lang

      this.user.lang = lang
      setUserLang(lang)

      if (lang == defaultLang._code) {
        Object.assign(i18n, defaultLang)
        return
      } else if (i18ns.has(lang)) {
        Object.assign(i18n, i18ns.get(lang))
        return
      } else {
        let strings = getI18nStrings(lang)
        if (null != strings && strings._version == defaultLang._version) {
          i18ns.set(lang, strings)
          Object.assign(i18n, strings)
          return
        }
        removeI18nStrings(lang)
      }

      axios({
        method: "GET",
        url: queryStatic("/static/strings/strings-" + lang + ".json"),
      }).then(function (resp) {
        i18ns.set(lang, resp.data)
        setI18nStrings(lang, resp.data)
        Object.assign(i18n, resp.data)
      }).catch(function (resp) {
        console.error("catch", resp)
        app.user.lang = old
        setUserLang(old)
      })
    },
  }
}).$mount('#app')