// app.js
// Created at 03/09/2021

const Home = { inject: ['i18n'], template: '<div>{{ i18n.thisIsHome }}</div>' }
const Foo = { inject: ['i18n'], template: '<div>{{ i18n.thisIsFoo }}</div>' }
const Bar = { inject: ['i18n'], template: '<div>{{ i18n.thisIsBar }}: {{ $route.params.id }}</div>' }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/my', name: 'my', component: Foo },
    { path: '/drafts/:id', name: 'drafts', component: Bar }
  ]
})

new Vue({
  provide: {
    i18n,
  },
  data: () => ({
    i18n,
  }),
  router,
  template: fgm_app,
  methods: {
    onToggle: function () {
      if (i18n.language === 'en') {
        Object.assign(i18n, zh)
      } else {
        Object.assign(i18n, en)
      }
    }
  }
}).$mount('#app')