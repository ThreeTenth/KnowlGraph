// vue_time.js

Vue.component('modal', {
  methods: {
    toggle() {
      this.$emit('toggle', false)
    },
  },

  template: com_modal,
})