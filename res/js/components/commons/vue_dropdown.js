// vue_dropdown.js

Vue.component('dropdown', {
  props: ['ignore'],

  data: function () {
    return {
      seen: false,
      pos: 'down',
    }
  },

  methods: {
    toggle() {
      if (this.seen) {
        return this.__hide()
      }
      return this.__show()
    },

    __show() {
      this.__update()
      this.seen = true
      setTimeout(() => document.addEventListener('click', this.__hide), 0);
    },

    __hide(e) {
      if (!e) return
      if (this.ignore) {
        var path = e.path
        var hide = true
        for (let index = 0; index < path.length; index++) {
          const element = path[index];
          if (element.className == "menu") {
            hide = false
            break
          }
        }
        if (hide) {
          this.seen = false
          document.removeEventListener('click', this.__hide);
        }
      } else {
        this.seen = false
        document.removeEventListener('click', this.__hide);
      }
    },

    __update() {
      let rect = this.$refs.dropdown.getBoundingClientRect()
      let top = rect.top
      let bottom = getBodyHeight() - rect.bottom
      if (top < 240 || 240 < bottom) {
        this.pos = 'down'
      } else {
        this.pos = 'up'
      }
    }
  },

  template: com_dropdown,
})