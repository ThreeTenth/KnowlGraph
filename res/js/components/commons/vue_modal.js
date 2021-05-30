// vue_modal.js

Vue.component('modal', {
  methods: {
    toggle() {
      this.$emit('toggle', false)
    },
  },

  template: com_modal,
})