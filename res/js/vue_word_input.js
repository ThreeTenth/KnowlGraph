// vue word input

const WordStateIdle = 0
const WordStateMayDel = 1
const WordStatePreDel = 2
const WordStateSafe = 3

Vue.component('words-input', {
  props: ['words'],

  data: function () {
    return {
      auto: false,
      values: this.words ? this.words : [],
      input: "",
      items: [],
      all: [],
      state: WordStateMayDel,
    }
  },

  created() { console.log(this.$props) },
  methods: {
    onWordChanged: function () {
      if ("" == this.input) {
        if (WordStateSafe == this.state) {
          this.state = WordStateIdle
        }
        this.auto = false
      } else {
        this.state = WordStateSafe
        if (2 <= this.input.length) {
          const items = []
          const regex = new RegExp('^' + this.word)
          this.all.forEach(item => {
            if (regex.test(item)) {
              items.push(item)
            }
          });
          this.items = items
          this.auto = true
        } else {
          this.auto = false
        }
      }
    },
    onRemoveWord: function (word) {
      const index = this.values.indexOf(word)
      if (index > -1) {
        this.values.splice(index, 1)
      }
      this.state = WordStateMayDel
    },
    onAdd: function () {
      var word = this.input.replace(/[\,\， ]*$/g, '')
      if ("" == word) {
        this.input = ""
        // todo Waring toast
        return
      }

      var words = this.values
      for (let index = 0; index < words.length; index++) {
        const element = words[index];
        if (element == word) {
          this.input = ""
          // todo Waring toast
          return
        }
      }

      this.values.push(word)
      this.input = ""
      this.state = WordStateMayDel
    },
    onDel: function () {
      switch (this.state) {
        case WordStateIdle:
          this.state = WordStateMayDel
          break;
        case WordStateMayDel:
          this.state = WordStatePreDel
          cancelDeleteWordTimeoutID = window.setTimeout(() => { this.state = WordStateMayDel }, 2000)
          break;
        case WordStatePreDel:
          const removeIndex = this.values.length
          this.values.splice(removeIndex - 1, 1)
          this.state = WordStateMayDel
          window.clearTimeout(cancelDeleteWordTimeoutID)
          break;
      }
    },
  },
  template: com_words_input,
})