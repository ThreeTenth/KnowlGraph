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

  // render(h) {
  //   const content = h('div', {
  //     staticClass: 'dialog__content',
  //   }, [
  //     this.$slots.default
  //   ])

  //   const dialog = h('div', {
  //     staticClass: 'modal-content',
  //   }, [
  //     this.$slots.default,
  //   ])

  //   const dialogContent = h('div', {
  //     staticClass: 'modal',
  //     ref: 'dialog-content',
  //     on: {
  //       click: (evt) => {
  //         if (!evt.target.closest('.modal-content') && !this.preventClose) {
  //           this.$emit('toggle', !this.seen)
  //           this.$emit('close')
  //         }
  //       }
  //     }
  //   }, [
  //     dialog
  //   ])

  //   const modalTransition = h('transition', {
  //     props: {
  //       name: 'vs-dialog'
  //     },
  //   }, [this.seen && dialogContent])

  //   return h('span', [modalTransition])
  // },
})