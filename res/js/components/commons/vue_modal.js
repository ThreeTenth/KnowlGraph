// vue_modal.js

Vue.component('modal', {
  model: {
    prop: 'seen',
    event: 'toggle'
  },

  props: {
    seen: Boolean
  },

  watch: {
    seen: function (val, oldVal) {
      if (val) {
        this.$nextTick(() => {
          const dialog = this.$refs['dialog-content']
          const target = document.body
          target.insertBefore(dialog, target.lastChild)
        })
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style['overflow-y'] = 'scroll'
      }
    },
  },

  beforeDestroy() {
    if (this.$el && this.$el.parentNode) {
      this.$el.parentNode.removeChild(this.$el)
    }
  },

  methods: {
    close() {
      this.$emit('toggle', false)
    },
  },

  template: com_modal,
})