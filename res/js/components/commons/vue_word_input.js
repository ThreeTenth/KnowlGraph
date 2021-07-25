// vue word input

const WordStateIdle = 0
const WordStateMayDel = 1
const WordStatePreDel = 2
const WordStateSafe = 3

Vue.component('words-input', {
  props: ['words', 'keywords'],

  data: function () {
    return {
      values: this.keywords ? this.keywords : [],
      input: "",
      newInput: "",
      items: [],
      deleteMark: null,
      __isDelProcess: false,
    }
  },

  methods: {
    onKeyup: function (event) {
      if ("Process" === event.key) {
        // 针对中文输入法的特殊处理
        // 在进入中文输入时，
        // vue 无法观测到正在输入的内容变化，
        // 但可以监听到按键变化。
        // 所以如果当前内容为空，
        // 但用户正在使用中文输入法输入内容，
        // 如果使用了删除键，
        // 这个删除事件就会被捕获，
        // 那么就会出现错误的删除操作。
        // 在 Vue 中，如果用户正在进行中文输入，
        // 那么每次按键事件都被前后触发两次，
        // 第一次是 key 始终为 "Process" 的按键事件，
        // 第二次是正常的按键事件，
        // 通过这个，可以过滤不需要的按键。
        // 在这里，就是删除键。

        this.$data.__isDelProcess = true
      }
    },
    onRemoveWord: function (word) {
      const index = this.values.indexOf(word)
      if (index > -1) {
        this.values.splice(index, 1)
      }
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
    },
    onDel: function () {
      if (this.$data.__isDelProcess) {
        this.$data.__isDelProcess = false
        return
      }
      if (0 < this.newInput.length) return
      if (null == this.deleteMark) {
        cancelDeleteWordTimeoutID = window.setTimeout(() => { this.deleteMark = null }, 2000)
        this.deleteMark = this.values.length - 1
      } else {
        this.values.splice(this.deleteMark, 1)
        window.clearTimeout(cancelDeleteWordTimeoutID)
        this.deleteMark = null
      }
    },
  },

  watch: {
    input: function (val, oldVal) {
      window.setTimeout(() => { this.newInput = val }, 100)
      if ("" == val) {
        this.items = []
      } else {
        this.deleteMark = null
        if (this.words) {
          if (2 <= val.length) {
            const items = []
            const regex = new RegExp('^' + val)
            this.words.forEach(word => {
              if (regex.test(word.name)) {
                items.push(word)
              }
            });
            this.items = items
         }
        }
      }
    },
  },

  template: com_words_input,
})