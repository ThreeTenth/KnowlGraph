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
        document.body.style.overflow = ''
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

  // template: com_modal,

  render(h) {
    const content = h('div', {
      staticClass: 'vs-dialog__content',
      // class: {
      //   notFooter: !this.$slots.footer
      // }
    }, [
      this.$slots.default
    ])

    const dialog = h('div', {
      staticClass: 'modal-content',
      // style: {
      //   width: this.width
      // },
      // class: {
      //   'vs-dialog--fullScreen': this.fullScreen,
      //   'vs-dialog--rebound': this.rebound,
      //   'vs-dialog--notPadding': this.notPadding,
      //   'vs-dialog--square': this.square,
      //   'vs-dialog--autoWidth': this.autoWidth,
      //   'vs-dialog--scroll': this.scroll,
      //   'vs-dialog--loading': this.loading,
      //   'vs-dialog--notCenter': this.notCenter,
      // }
    }, [
      // this.loading && loading,
      // !this.notClose && close,
      // this.$slots.header && header,
      this.$slots.default,
      // this.$slots.footer && footer
    ])

    const dialogContent = h('div', {
      staticClass: 'modal',
      ref: 'dialog-content',
      // class: {
      //   blur: this.blur,
      //   fullScreen: this.fullScreen,
      // },
      on: {
        click: (evt) => {
          if (!evt.target.closest('.modal-content') && !this.preventClose) {
            this.$emit('toggle', !this.seen)
            this.$emit('close')
          }
        }
      }
    }, [
      dialog
    ])

    const modalTransition = h('transition', {
      props: {
        name: 'vs-dialog'
      },
    }, [this.seen && dialogContent])

    return h('span', [modalTransition])
  },
})