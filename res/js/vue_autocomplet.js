// vue_time.js

Vue.component('autocomplet', {
  props: ["placeholder", "source"],

  data: function () {
    return {
      input: "",
      items: [],
      auto: false,
      pos: 'down',
    }
  },

  mounted() {
    this.__update()
  },

  activated() {
    this.__update()
  },

  methods: {
    onSelect(item) {
      this.$emit('item', item, val => {
        this.input = val
      })
      this.auto = false
      this.$emit('select', item)
    },

    onKeyup() {
      this.auto = "" != this.input

      // if (!this.auto) return

      const items = []
      const regex = new RegExp('^' + this.input)
      this.source.forEach(item => {
        let value = item
        this.$emit('item', item, val => {
          value = val
        })

        if (regex.test(value)) {
          items.push(item)
        }
      });
      this.items = items
      this.auto = 0 != items.length

      this.__update()
    },

    __update() {
      let len = this.items.length

      if (0 == len ) return

      let autocomplete = this.$refs.autocomplete
      let top = autocomplete.getBoundingClientRect().top
      let bottom = autocomplete.getBoundingClientRect().bottom

      let preHeight = len * 30

      console.log(top, bottom, preHeight)

      if (preHeight < bottom) {
        this.pos = "down"
      } else {
        this.pos = "up"
      }
    },
  },

  template: com_autocomplet,
})