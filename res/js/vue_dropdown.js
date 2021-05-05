// vue_time.js

Vue.component('dropdown', {
  props: ['ignore'],

  data: function () {
    return {
      seen: false,
      pos: 'down',
    }
  },
  computed: {
  },
  mounted() {
    this.__update()
  },
  activated() {
    this.__update()
  },
  methods: {
    toggle() {
      if (this.seen) {
        return this.__hide()
      }
      return this.__show()
    },

    __show() {
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
          if (element.className == "drop-content") {
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
      let top = this.$refs.dropdown.getBoundingClientRect().top
      if (top < 240) {
        this.pos = 'down'
      } else {
        this.pos = 'up'
      }
    }
  },

  template: com_dropdown,
})