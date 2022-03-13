// vue_new_node.js

Vue.component('new-node', {
  props: ['article'],

  data: function () {
    return {
      path: null,
      level: -1,
      value: "",
      isNewNode: false,
    }
  },

  computed: {
    isRoot: function () {
      return this.level == -1 || this.level == 0
    },
  },

  methods: {
    onOpenNode(index, node) {
      var path = this.path
      var level = this.level

      let nexts = node.edges.nexts
      if (!nexts) {
        node.loading = true
        getNodeNexts(node.id).
          then(resp => this.loadNodeNexts2Path(resp, node)).
          catch(this.loadError)
      }

      path.splice(index + 1)
      path.push(node)
      level = index + 1

      this.path = path
      this.level = level
      this.isNewNode = false
    },

    onSelectWord(word) {
      this.value = word
    },

    onBackPrev(i = -1) {
      this.isNewNode = false
      if (-1 == i) {
        this.level--
        this.path.pop()
      } else {
        this.level = i
        this.path.splice(i + 1)
      }
    },

    onNewNode() {
      this.isNewNode = true
    },

    getWordName(word, callback) {
      callback(word.name)
    },

    onCancelNewNode() {
      this.isNewNode = false
    },

    onSubmitNewNode() {
      if (!this.value) return;

      let node = this.path[this.level]
      axios({
        method: "PUT",
        url: queryRestful("/v1/node"),
        data: {
          wordId: this.value.id,
          nodeId: node.id,
        },
      }).then(resp => this.addNodeSuccessed(resp, node)).
        catch(this.addNodeFailured)
    },

    addNode(root, node) {
      let nexts = root.edges.nexts

      if (nexts) {
        nexts.push(node)
      } else {
        nexts = [node]
        root.edges.nexts = nexts
      }
    },

    addNodeSuccessed(resp, root) {
      let archive = resp.data
      let node = archive.edges.node
      node.edges.word = this.value
      this.addNode(root, node)
      this.isNewNode = false
    },
    addNodeFailured(resp) {
      this.toast("Error(" + resp.code + "): " + resp.data, "error")
      console.error(resp);
    },

    reload() {
      var root = {
        loading: true,
        edges: {
          word: {
            name: "已选择"
          }
        }
      }
      var path = [root]
      var level = 0

      this.path = path
      this.level = level
      getNodeNexts(0).
        then(resp => this.loadNodeNexts2Path(resp, root)).
        catch(this.loadError)
    },

    loadNodeNexts(resp, root) {
      let node = resp.data

      root.loading = false
      root.edges.nexts = node.edges.nexts
    },

    loadNodeNexts2Path(resp, root) {
      this.loadNodeNexts(resp, root)

      let path = this.path

      path.pop()
      path.push(root)

      this.path = path
    },

    loadError(err) {
      this.toast(err, "error")
      console.error(err);
    },

    onArchive() {
      if (this.isRoot) {
        this.toast("请选择一个节点", "error")
        return
      }
      if (!this.article) {
        this.toast("没有指定文章归档")
        return
      }

      var node = this.path[this.level]
      var article = this.article

      putNodeArticle(node.id, article.id).
        then(resp => this.archiveSuccess(resp, node)).
        catch(this.archiveFailure)
    },

    archiveSuccess(resp, node) {
      this.$emit('successed', node)
    },
    archiveFailure(err) {
      this.$emit('failured', err)
    },
  },

  created() {
    this.reload()
    setTimeout(() => {
      this.$el.addEventListener("click", this.onCancelNewNode, { passive: false })
    }, 0);
  },

  beforeDestroy() {
    this.$el.removeEventListener("click", this.onCancelNewNode)
  },

  template: com_new_node,
})