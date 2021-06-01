// vue_autocomplet.js

Vue.component('autocomplet', {
  props: ["placeholder", "source"],

  data: function () {
    return {
      input: "",
      items: [],
      auto: false,
      selected: 0,
      pos: 'down',
      watchTime: 0,
    }
  },

  watch: {
    items: function (val, oldVal) {
      let len = val.length

      if (0 == len) return

      this.$refs.selection.style.maxHeight = ''
      this.pos = "down"

      this.$nextTick(() => {
        let selection = this.$refs.selection
        let rect = selection.getBoundingClientRect()
        let windowHeight = document.documentElement.clientHeight
        let top = rect.top
        let preHeight = rect.height

        if (top + preHeight < windowHeight) {
          this.pos = "down"
        } else {
          this.pos = "up"

          selection.style.maxHeight = rect.top + 'px'
        }
      })
    },
  },

  mounted() {
    // this.__update()
  },

  activated() {
    // this.__update()
  },

  methods: {
    onSelect(item) {
      this.$emit('item', item, val => {
        this.input = val
      })
      this.$emit('select', item)
      this.items = []
      this.auto = false
    },

    onKeyTab() {
      console.log("on tab")
    },

    onHotKey(e) {
      switch (e.keyCode) {
        case 38:  // up
          if (0 == this.items.length) return
          if (0 < this.selected) {
            this.selected -= 1
          } else {
            this.selected = this.items.length - 1
          }
          this.$refs.selection.children[this.selected].scrollIntoView()
          e.preventDefault()
          break;
        case 9:   // tab
        case 40:  // down
          if (0 == this.items.length) return
          if (this.items.length <= this.selected + 1) {
            this.selected = 0
          } else {
            this.selected += 1
          }
          this.$refs.selection.children[this.selected].scrollIntoView()
          e.preventDefault()
          break;
        case 13:  // enter
          if (0 == this.items.length) return
          this.onSelect(this.items[this.selected])
          e.preventDefault()
          break;
      }
    },

    onChanged(e) {
      // console.log("onChanged", e.code, e.key, e.keyCode, e.target.value)

      switch (e.keyCode) {
        case 38:  // up
        case 9:   // tab
        case 40:  // down
        case 13:  // enter
          return
      }

      const items = []
      const regex = new RegExp('^' + e.target.value)
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
      this.selected = 0
      this.auto = true
      // this.__update()
    },

    // __update() {
    //   let len = this.items.length

    //   if (0 == len) return

    //   let autocomplete = this.$refs.autocomplete
    //   let bottom = autocomplete.getBoundingClientRect().bottom

    //   let preHeight = len * 30

    //   if (preHeight < bottom) {
    //     this.pos = "down"
    //   } else {
    //     this.pos = "up"
    //   }
    // },
  },

  template: com_autocomplet,
})