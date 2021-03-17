// vue tag input

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

Vue.component('tags-input', {
  props: ['tags'],

  data: function () {
    return {
      auto: false,
      values: this.tags,
      input: "",
      items: [],
      all: [],
      state: TagStateMayDel,
    }
  },

  created() { console.log(this.$props) },
  methods: {
    onTagChanged: function () {
      if ("" == this.input) {
        if (TagStateSafe == this.state) {
          this.state = TagStateIdle
        }
        this.auto = false
      } else {
        this.state = TagStateSafe
        if (2 <= this.input.length) {
          const items = []
          const regex = new RegExp('^' + this.tag)
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
    onRemoveTag: function (tag) {
      const index = this.values.indexOf(tag)
      if (index > -1) {
        this.values.splice(index, 1)
      }
      this.state = TagStateMayDel
    },
    onAdd: function () {
      var tag = this.input.replace(/[\,\， ]*$/g, '')
      if ("" == tag) {
        this.input = ""
        // todo Waring toast
        return
      }

      var tags = this.values
      for (let index = 0; index < tags.length; index++) {
        const element = tags[index];
        if (element == tag) {
          this.input = ""
          // todo Waring toast
          return
        }
      }

      this.values.push(tag)
      this.input = ""
      this.state = TagStateMayDel
    },
    onDel: function () {
      switch (this.state) {
        case TagStateIdle:
          this.state = TagStateMayDel
          break;
        case TagStateMayDel:
          this.state = TagStatePreDel
          cancelDeleteTagTimeoutID = window.setTimeout(() => { this.state = TagStateMayDel }, 2000)
          break;
        case TagStatePreDel:
          const removeIndex = this.values.length
          this.values.splice(removeIndex - 1, 1)
          this.state = TagStateMayDel
          window.clearTimeout(cancelDeleteTagTimeoutID)
          break;
      }
    },
  },
  template: com_tags_input,
})