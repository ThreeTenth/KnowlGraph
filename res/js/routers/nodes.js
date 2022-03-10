// nodes.js

const Nodes = {
  props: {
    id: {
      type: Number,
      default: 0,
    }
  },

  data: function () {
    return {
      pageStatus: 0,
      node: null,
    }
  },

  watch: {
    id(value, oldVal) {
      this.reload()
    },
  },

  methods: {
    reload() {
      getNodeNexts(this.id).then(this.loadData).catch(this.loadError)
    },

    loadData(response) {
      this.pageStatus = response.status
      this.node = response.data
    },
    loadError(err) {
      this.pageStatus = getStatus4Error(err)
      console.error(err);
    },
  },
  created() {
    this.reload()
  },
  template: fgm_nodes,
}