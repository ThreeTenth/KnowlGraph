// vote.js

const Vote = {
  data: function () {
    return {
    }
  },
  computed: {
    version: function () {
      var ver = this.vote.value.edges.version
      this.languages.forEach(element => {
        if (element.code == ver.lang) {
          ver.lang = element.name
        }
      });
      return ver
    },
  },
  methods: {
    onAllow() {
      this.__postVote("allowed")
    },
    onRejecte() {
      this.__postVote("rejected")
    },
    onAbstained() {
      this.__postVote("abstained")
    },
    __postVote(status) {
      if (!this.version.exist) return

      axios({
        method: "POST",
        url: queryRestful("/v1/vote"),
        data: {
          id: this.vote.value.id,
          status: status,
        },
      }).then((resp) => {
        Object.assign(_voteObservable, { exist: false, value: null })
        this.back()
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },
  created() {
    document.title = "表决空间 -- KnowlGraph"
    // document.body.style.overflow = 'hidden'
  },

  beforeDestroy() {
    // document.body.style['overflow-y'] = 'scroll'
  },
  template: fgm_vote,
}