// vue_modal.js

Vue.component('modal', {
  model: {
    prop: 'seen',
    event: 'toggle'
  },

  props: {
    seenClose: {
      type: Boolean,
      default: true,
    },
    clazz: String,
    seen: Boolean
  },

  watch: {
    seen: function (val, oldVal) {
      if (val) {
        this.$nextTick(() => {
          const dialog = this.$refs['dialog-content']
          const target = document.body
          target.insertBefore(dialog, target.lastChild)
          dialog.addEventListener("click", this.clickDialog, { passive: false })
        })
        document.body.style.overflow = 'hidden'
      }
    },
  },

  beforeDestroy() {
    if (this.$el && this.$el.parentNode) {
      this.$el.parentNode.removeChild(this.$el)
    }
    document.body.style['overflow-y'] = 'scroll'
  },

  methods: {
    close() {
      if (this.seenClose) {
        const dialog = this.$refs['dialog-content']
        dialog.removeEventListener("click", this.clickDialog)
        document.body.style['overflow-y'] = 'scroll'

        this.$emit('toggle', false)
      }
    },
    clickDialog(e) {
      const content = this.$refs['dialog-content']
      if (content != e.target) {
        return
      }

      e.stopPropagation()
      e.preventDefault()
      this.close()
    },
  },

  template: com_modal,
})