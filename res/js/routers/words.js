// words.js

const Words = {
  computed: {
    publicWords() {
      if (!this.user.words) return null
      var pws = []
      this.user.words.forEach(word => {
        if (word.status == 'public')
          pws.push(word)
      })
      return pws.length == 0 ? null : pws
    },
    yourPrivateWords() {
      if (!this.user.words) return null
      var pws = []
      this.user.words.forEach(word => {
        if (word.status == 'private')
          pws.push(word)
      })
      return pws.length == 0 ? null : pws
    },
  },
  template: fgm_words,
}