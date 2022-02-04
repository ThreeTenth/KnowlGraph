// vote.js

const Vote = {
  data: function () {
    return {
      codeofConduct: null,
      overruleSelected: [],
      allowDropdown: false,
      overruleDropdown: false,
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
      this.allowDropdown = !this.allowDropdown
    },
    onOverrule() {
      this.overruleDropdown = !this.overruleDropdown
      this.__postVote("overruled", this.overruleSelected)
    },
    onAbstain() {
      this.__postVote("abstained")
    },
    __postVote(status, code = []) {
      postVote(this.vote.value.id, status, code, () => {
        Object.assign(_voteObservable, { exist: false, value: null })
        this.back()
      }, (err) => {
        console.log(err)
      })
    },
  },
  created() {
    document.title = "表决空间 -- KnowlGraph"

    getCovenant("zh", (resp) => { this.codeofConduct = resp.data }, (err) => { console.error(err) })
  },

  beforeDestroy() {
  },
  template: fgm_vote,
}