// vue tag input

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

var en = {
  language: "en",
  addATag: "Add a tag",
  thisIsHome: "This is Home",
  thisIsFoo: "This is Foo",
  thisIsBar: "This is Bar",
  namedRoutes: "Named Routes",
  currentRouteName: "Current route name",
  home: "Home",
  drafts: "Drafts",
  my: "Bar",
}

var zh = {
  language: "zh",
  addATag: "添加一个标签",
  thisIsHome: "这是首页",
  thisIsFoo: "这是Foo页",
  thisIsBar: "这是Bar页",
  namedRoutes: "路由名称",
  currentRouteName: "当前路由名称",
  home: "首页",
  drafts: "草稿箱",
  my: "我的",
}

const i18n = Vue.observable({
	...en,
})

Vue.component('tags-input', {
  inject: ['i18n'],
  data: function () {
    return {
      auto: false,
      values: [],
      input: "",
      items: [],
      all: [],
      state: TagStateMayDel,
    }
  },
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