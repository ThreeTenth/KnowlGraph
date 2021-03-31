// vue_time.js

Vue.component('dropdown', {
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

    __hide() {
      this.seen = false
      document.removeEventListener('click', this.__hide);
    },

    __update() {
      let top = this.$refs.dropdown.getBoundingClientRect().top
      console.log(top)
      if (top < 240) {
        this.pos = 'down'
      } else {
        this.pos = 'up'
      }
    }
  },

  template: com_dropdown,
})