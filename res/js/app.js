// app.js
// Created at 03/09/2021

const Home = { inject: ['i18n'], template: fgm_home }
const My = { inject: ['i18n'], template: fgm_my }
const Drafts = { inject: ['i18n'], template: fgm_user_drafts }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/my', name: 'my', component: My },
    { path: '/drafts/:id', name: 'drafts', component: Drafts }
  ]
})

var app = new Vue({
  provide: {
    i18n,
  },
  data: {
    user: {
      lang: getUserLang(),
    },
    languages: languages,
    i18n,
  },
  router,
  template: fgm_app,
  methods: {
    onToggle: function () {
      if (i18n.language === 'en') {
        Object.assign(i18n, zh)
      } else {
        Object.assign(i18n, en)
      }
    },
    onSelectUserLang: function (event) {
      const lang = event.target.value
      const old = this.user.lang

      // todo need cache

      axios({
        method: "GET",
        url: queryStatic("/static/strings/strings-" + lang + ".json"),
      }).then(function (resp) {
        Object.assign(i18n, resp.data)
      }).catch(function (resp) {
        console.error("catch", resp)
        app.user.lang = old
      })

      this.user.lang = lang
    },
  }
}).$mount('#app')