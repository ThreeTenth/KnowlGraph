// vue_overlay.js

Vue.component('overlay', {

  created() {
    document.body.style.overflow = 'hidden'
  },

  beforeDestroy() {
    document.body.style['overflow-y'] = 'scroll'
  },

  template: com_overlay,
})