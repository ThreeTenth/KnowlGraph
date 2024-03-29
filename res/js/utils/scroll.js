
class ElementIndex {
  constructor(index, overflow, offset, start) {
    this.index = index
    this.overflow = overflow
    this.offset = offset
    this.start = start
  }
}

const ScrollYOffset = 3

function indexStartAtWindow(elements) {
  let index = 0
  for (; index < elements.length; index++) {
    if (elements[index].offsetTop + elements[index].offsetHeight <= window.scrollY + ScrollYOffset) continue
    if (elements[index].offsetTop <= window.scrollY + ScrollYOffset) break
  }

  index = elements.length <= index ? elements.length - 1 : index
  let offset = (elements[index].offsetTop + elements[index].offsetHeight - window.scrollY - window.getBodyHeight() - ScrollYOffset)
  let start = (window.scrollY - elements[index].offsetTop)
  return new ElementIndex(index, 0 < offset, offset, start)
}

function indexEndAtWindow(elements) {
  let index = elements.length - 1
  for (; 0 <= index; index--) {
    if (elements[index].offsetTop + ScrollYOffset < window.scrollY + window.getBodyHeight()) break
  }

  index = index < 0 ? 0 : index
  let offset = (window.scrollY - elements[index].offsetTop - ScrollYOffset)
  return new ElementIndex(index, 0 < offset, offset, offset + ScrollYOffset)
}