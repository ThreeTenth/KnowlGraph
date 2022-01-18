// vue_fullpage.js

Vue.component('fullpage', {

  data: function () {
    return {
      elements: [],
    }
  },

  methods: {
    downMove(deltaY, behavior = "smooth") {
      let startElement = indexStartAtWindow(this.elements)
      if (startElement.overflow) {
        if (deltaY < startElement.offset) {
          window.scrollBy({ top: deltaY, behavior: behavior })
        } else {
          window.scrollBy({ top: startElement.offset + ScrollYOffset, behavior: behavior })
        }
      } else {
        let pageIndex = getValueInRange(startElement.index + 1, 0, this.elements.length)
        this.elements[pageIndex].scrollIntoView({ behavior: "smooth" })
      }
    },

    upMove(deltaY, behavior = "smooth") {
      let endElement = indexEndAtWindow(this.elements)
      if (endElement.overflow) {
        if (-endElement.offset < deltaY) {
          window.scrollBy({ top: deltaY, behavior: behavior })
        } else {
          window.scrollBy({ top: -endElement.offset - ScrollYOffset, behavior: behavior })
        }
      } else {
        let pageIndex = getValueInRange(endElement.index - 1, 0, this.elements.length)
        this.elements[pageIndex].scrollIntoView({ behavior: "smooth" })
      }
    },

    keydownHandle(e) {
      if (e.altKey && (e.keyCode == 37 || e.keyCode == 39)) {
        // 浏览器前进后退的默认行为
        return
      }
      switch (e.keyCode) {
        case 9:  // Tab
        case 35:  // End
        case 36:  // Home
          // 禁止 Tab, Home 和 End 键的使用。
          // 此功能会导致阅读不完整。
          e.preventDefault()
          break
        case 33:  // 上页
          e.preventDefault()
          this.upMove(-(window.getBodyHeight()))
          break
        case 34:  // 下页
          e.preventDefault()
          this.downMove(window.getBodyHeight())
          break
        case 37:  // 左
        case 38:  // 上
          e.preventDefault()
          this.upMove(-200)
          break
        case 32:  // 空格键
        case 39:  // 右
        case 40:  // 下
          e.preventDefault()
          this.downMove(200)
          break

        default:
          return
      }
    },
  },

  created() {
    document.addEventListener('keydown', this.keydownHandle, { passive: false })    //passive 参数不能省略，用来兼容ios和android

    this.$nextTick(() => {
      let slots = this.$slots["default"]
      for (let index = 0; index < slots.length; index++) {
        const element = slots[index];
        if (element.tag) {
          this.elements.push(element.elm)
        }
      }
    })
  },

  beforeDestroy() {
    document.removeEventListener('keydown', this.keydownHandle)
  },

  template: com_fullpage,
})