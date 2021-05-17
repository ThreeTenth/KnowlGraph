// vue_time.js

Vue.component('autocomplet', {
  props: ["placeholder", "source"],

  data: function () {
    return {
      input: "",
      items: [],
      auto: false,
    }
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

      if (!this.auto) return

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
    }
  },

  template: com_autocomplet,
})