// app.js
// Created at 03/09/2021

const Home = { template: '<div>This is Home</div>' }
const Foo = { template: '<div>This is Foo</div>' }
const Bar = { template: '<div>This is Bar {{ $route.params.id }}</div>' }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/my', name: 'foo', component: Foo },
    { path: '/drafts/:id', name: 'bar', component: Bar }
  ]
})

new Vue({
  router,
  template: fgm_app
}).$mount('#app')