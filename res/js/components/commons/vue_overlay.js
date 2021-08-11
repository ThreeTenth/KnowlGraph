// vue_overlay.js

Vue.component('overlay', {

  data: function () {
    return {
      body_overflow_y: 'scroll',
    }
  },

  created() {
    if (document.body.style['overflow-y']) {
      this.body_overflow_y = document.body.style['overflow-y']
    }
    document.body.style.overflow = 'hidden'
  },

  beforeDestroy() {
    document.body.style['overflow-y'] = this.body_overflow_y
  },

  template: com_overlay,
})