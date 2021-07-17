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
      const _this = this;
      if (!this.selected || !this.selected.edges || !this.selected.edges.node) return;
      if (!this.value) return;

      axios({
        method: "PUT",
        url: queryRestful("/v1/node"),
        data: {
          wordId: this.value.id,
          nodeId: this.selected.edges.node.id,
        },
      }).then(function (resp) {
        _this.$vs.notification({
          color: 'success',
          position: 'bottom-right',
          title: "Success",
        })
        _this.$emit('successed')
      }).catch(function (resp) {
        _this.$vs.notification({
          color: 'danger',
          position: 'bottom-right',
          title: "Failure",
          text: resp.data
        })
        _this.$emit('failured')
      })
    },
  },

  template: com_new_node,
})