// vue_dropdown.js

Vue.component('dropdown', {
  props: ["trigger"],

  data: function () {
    return {
      seen: false,
      pos: 'down',
    }
  },

  watch: {
    trigger: function (val, oldVal) {
      if (val == oldVal) return
      this.seen = !this.seen
      if (!this.seen) {
        document.removeEventListener('click', this.toggle);
      }
    },
  },

  methods: {
    toggle() {
      if (this.seen) {
        this.__hide()
      } else {
        this.__show()
      }
    },

    __show() {
      this.__update()
      this.seen = true
      setTimeout(() => document.addEventListener('click', this.toggle), 0);
    },

    __hide() {
      this.seen = false
      document.removeEventListener('click', this.toggle);
    },

    __update() {
      let rect = this.$el.getBoundingClientRect()
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