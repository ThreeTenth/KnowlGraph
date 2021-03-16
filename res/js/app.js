// app.js
// Created at 03/09/2021

const plugin = {
  install: function (Vue, options) {
    Vue.prototype.format = function (format, val) {
      let now = new Date()
      let date1 = new Date(timestamp)
      let timer = (now - date1) / 1000
      console.log(timer)
      let tip = ''

      if (timer <= 0) {
        tip = '刚刚'
      } else if (Math.floor(timer / 60) <= 0) {
        tip = '刚刚'
      } else if (Math.floor(timer < 3600)) {
        tip = Math.floor(timer / 60) + '分钟前'
      } else if (timer >= 3600 && timer < 86400) {
        tip = Math.floor(timer / 3600) + '小时前'
      } else if (timer / 86400 <= 31) {
        tip = Math.floor(timer / 86400) + '天前'
      } else {
        tip = "too lang ago"
      }

      return tip
    }

    Vue.prototype.i18n = i18n
  }
}

Vue.use(plugin)

const Index = { template: logined ? fgm_home : fgm_index }
const Home = { template: fgm_home }
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
    { path: '/d/:id/edit', name: 'editDraft', component: EditDraft, props: true },
    { path: '/d/:id/history', name: 'draftHistories', component: DraftHistories, props: true },
    { path: '/d/:id/history/:hid', name: 'draftHistory', component: DraftHistory, props: true },
    { path: '/d/:id/publish', name: 'publishArticle', component: PublishArticle },
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