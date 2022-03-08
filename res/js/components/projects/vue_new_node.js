// vue_new_node.js

Vue.component('new-node', {
  data: function () {
    return {
      selected: null,
      value: "",
    }
  },

  methods: {
    onSelectPrev(archive) {
      this.selected = this.selected == archive ? null : archive
    },

    onSelectWord(word) {
      this.value = word
    },

    getWordName(word, callback) {
      callback(word.name)
    },

    onSubmit() {
      if (!this.value) return;

      var body = {
        wordId: this.value.id,
      }

      if (this.selected && this.selected.edges && this.selected.edges.node) {
        body.nodeId = this.selected.edges.node.id
      }

      axios({
        method: "PUT",
        url: queryRestful("/v1/node"),
        data: body,
      }).then((resp) => {
        this.toast("Successed")
        this.$emit('successed', resp.data)
      }).catch((resp) => {
        this.toast("Failured")
        this.$emit('failured', resp.code, resp.data)
      })
    },
  },

  template: com_new_node,
})