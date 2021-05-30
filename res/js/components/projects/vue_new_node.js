// vue_new_node.js

Vue.component('new-node', {
  data: function () {
    return {
      selected: null,
      value: "",
    }
  },

  methods: {
    onSelectPrev: function (archive) {
      this.selected = this.selected == archive ? null : archive
    },

    onSelectWord(word) {
      this.value = word
    },

    getWordName(word, callback) {
      callback(word.name)
    }
  },

  template: com_new_node,
})