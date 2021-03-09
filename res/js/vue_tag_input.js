// vue tag input

const TagStateIdle = 0
const TagStateMayDel = 1
const TagStatePreDel = 2
const TagStateSafe = 3

function tagChanged() {
  if ("" == app.publish.tags.new) {
    if (TagStateSafe == app.publish.tags.state) {
      app.publish.tags.state = TagStateIdle
    }
    app.publish.tags.auto = false
  } else {
    app.publish.tags.state = TagStateSafe
    if (2 <= app.publish.tags.new.length) {
      const items = []
      const regex = new RegExp('^' + app.publish.tags.tag)
      app.publish.tags.all.forEach(item => {
        if (regex.test(item)) {
          items.push(item)
        }
      });
      app.publish.tags.items = items
      app.publish.tags.auto = true
    } else {
      app.publish.tags.auto = false
    }
  }
}

function removeTag(tag) {
  const index = app.publish.tags.values.indexOf(tag)
  if (index > -1) {
    app.publish.tags.values.splice(index, 1)
  }
  app.publish.tags.state = TagStateMayDel
}

function addTag() {
  // var tag = app.publish.tags.new.trim()                // ",vue,,,,,   ,web, golang      ,,  ,,rest,"
  // tag = tag.replace(/\, +,/g, '')         // ",vue,,,,web, golang      ,,rest,"
  // tag = tag.replace(/\,+/g, ',')          // ",vue,web, golang      ,rest,"
  // tag = tag.replace(/^\,+|\,+$/g, '')     // "vue,web, golang      ,rest"
  // tag = tag.replace(/ *\, */g, ',')       // "vue,web,golang,rest"
  //                                         // "gin    ", "gin  ,", "gin,", "gin,， ，  ,"
  // tag = tag.replace(/ *\,*$/g, '')        // "gin"

  var tag = app.publish.tags.new.replace(/[\,\， ]*$/g, '')
  if ("" == tag) {
    app.publish.tags.new = ""
    // todo Waring toast
    return
  }

  var tags = app.publish.tags.values
  for (let index = 0; index < tags.length; index++) {
    const element = tags[index];
    if (element == tag) {
      app.publish.tags.new = ""
      // todo Waring toast
      return
    }
  }

  app.publish.tags.values.push(tag)
  app.publish.tags.new = ""
  app.publish.tags.state = TagStateMayDel
}

function deleteTag() {
  switch (app.publish.tags.state) {
    case TagStateIdle:
      app.publish.tags.state = TagStateMayDel
      break;
    case TagStateMayDel:
      app.publish.tags.state = TagStatePreDel
      cancelDeleteTagTimeoutID = window.setTimeout(cancelDeleteTag, 2000)
      break;
    case TagStatePreDel:
      const removeIndex = app.publish.tags.values.length
      app.publish.tags.values.splice(removeIndex - 1, 1)
      app.publish.tags.state = TagStateMayDel
      window.clearTimeout(cancelDeleteTagTimeoutID)
      break;
  }
}

function cancelDeleteTag() {
  app.publish.tags.state = TagStateMayDel
}