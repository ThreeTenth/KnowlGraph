// vue_time.js

Vue.component('new-node', {
  data: function () {
    return {
      selected: 0,
      value: "",
    }
  },

  methods: {
    onSelectPrev: function (archive) {
      this.selected = this.selected == archive.id ? 0 : archive.id
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