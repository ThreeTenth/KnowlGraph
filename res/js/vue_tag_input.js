// vue tag input

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

var tag_layout = new Vue({
  el: '#tag_layout',
  data: {
    seen: false,
    auto: false,
    tags: [],
    tag: "",
    items: [],
    all: [],
    state: TagStateMayDel,
  },
  methods: {
    onTagChanged: function () {
      tagChanged()
    },
    onRemoveTag: function (tag) {
      removeTag(tag)
    },
    onAdd: function () {
      addTag()
    },
    onDel: function () {
      deleteTag()
    },
  },
})

function tagChanged() {
  if ("" == tag_layout.tag) {
    if (TagStateSafe == tag_layout.state) {
      tag_layout.state = TagStateIdle
    }
    tag_layout.auto = false
  } else {
    tag_layout.state = TagStateSafe
    if (2 <= tag_layout.tag.length) {
      const items = []
      const regex = new RegExp('^' + tag_layout.tag)
      tag_layout.all.forEach(item => {
        if (regex.test(item)) {
          items.push(item)
        }
      });
      tag_layout.items = items
      tag_layout.auto = true
    } else {
      tag_layout.auto = false
    }
  }
}

function removeTag(tag) {
  const index = tag_layout.tags.indexOf(tag)
  if (index > -1) {
    tag_layout.tags.splice(index, 1)
  }
  tag_layout.state = TagStateMayDel

  postBlur()
}

function addTag() {
  // var tag = tag_layout.tag.trim()         // ",vue,,,,,   ,web, golang      ,,  ,,rest,"
  // tag = tag.replace(/\, +,/g, '')         // ",vue,,,,web, golang      ,,rest,"
  // tag = tag.replace(/\,+/g, ',')          // ",vue,web, golang      ,rest,"
  // tag = tag.replace(/^\,+|\,+$/g, '')     // "vue,web, golang      ,rest"
  // tag = tag.replace(/ *\, */g, ',')       // "vue,web,golang,rest"
  //                                         // "gin    ", "gin  ,", "gin,", "gin,， ，  ,"
  // tag = tag.replace(/ *\,*$/g, '')        // "gin"

  var tag = tag_layout.tag.replace(/[\,\， ]*$/g, '')
  if ("" == tag) {
    tag_layout.tag = ""
    return
  }

  tag_layout.tags.push(tag)
  tag_layout.tag = ""
  tag_layout.state = TagStateMayDel

  postBlur()
}

function deleteTag() {
  switch (tag_layout.state) {
    case TagStateIdle:
      tag_layout.state = TagStateMayDel
      break;
    case TagStateMayDel:
      tag_layout.state = TagStatePreDel
      cancelDeleteTagTimeoutID = window.setTimeout(cancelDeleteTag, 2000)
      break;
    case TagStatePreDel:
      const removeIndex = tag_layout.tags.length
      tag_layout.tags.splice(removeIndex - 1, 1)
      tag_layout.state = TagStateMayDel
      window.clearTimeout(cancelDeleteTagTimeoutID)

      postBlur()
      break;
  }
}

function cancelDeleteTag() {
  tag_layout.state = TagStateMayDel
}