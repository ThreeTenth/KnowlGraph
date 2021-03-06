// vue tag input

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

function tagChanged() {
  if ("" == app.tag) {
    if (TagStateSafe == app.state) {
      app.state = TagStateIdle
    }
    app.auto = false
  } else {
    app.state = TagStateSafe
    if (2 <= app.tag.length) {
      const items = []
      const regex = new RegExp('^' + app.tag)
      app.all.forEach(item => {
        if (regex.test(item)) {
          items.push(item)
        }
      });
      app.items = items
      app.auto = true
    } else {
      app.auto = false
    }
  }
}

function removeTag(tag) {
  const index = app.tags.indexOf(tag)
  if (index > -1) {
    app.tags.splice(index, 1)
  }
  app.state = TagStateMayDel

  postBlur()
}

function addTag() {
  // var tag = app.tag.trim()                // ",vue,,,,,   ,web, golang      ,,  ,,rest,"
  // tag = tag.replace(/\, +,/g, '')         // ",vue,,,,web, golang      ,,rest,"
  // tag = tag.replace(/\,+/g, ',')          // ",vue,web, golang      ,rest,"
  // tag = tag.replace(/^\,+|\,+$/g, '')     // "vue,web, golang      ,rest"
  // tag = tag.replace(/ *\, */g, ',')       // "vue,web,golang,rest"
  //                                         // "gin    ", "gin  ,", "gin,", "gin,， ，  ,"
  // tag = tag.replace(/ *\,*$/g, '')        // "gin"

  var tag = app.tag.replace(/[\,\， ]*$/g, '')
  if ("" == tag) {
    app.tag = ""
    return
  }

  app.tags.push(tag)
  app.tag = ""
  app.state = TagStateMayDel

  postBlur()
}

function deleteTag() {
  switch (app.state) {
    case TagStateIdle:
      app.state = TagStateMayDel
      break;
    case TagStateMayDel:
      app.state = TagStatePreDel
      cancelDeleteTagTimeoutID = window.setTimeout(cancelDeleteTag, 2000)
      break;
    case TagStatePreDel:
      const removeIndex = app.tags.length
      app.tags.splice(removeIndex - 1, 1)
      app.state = TagStateMayDel
      window.clearTimeout(cancelDeleteTagTimeoutID)

      postBlur()
      break;
  }
}

function cancelDeleteTag() {
  app.state = TagStateMayDel
}