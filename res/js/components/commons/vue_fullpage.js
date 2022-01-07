// vue_fullpage.js

Vue.component('fullpage', {

  data: function () {
    return {
      scrollTimer: null,
      scrollTime: new Date(),
      elements: [],
      isIntercept: false,
    }
  },

  methods: {
    scrollStart(e) {
    },
    scrollDo(e) {
      if (this.isIntercept) return
      if (0 < e.deltaY) {
        // down

        let startElement = indexStartAtWindow(this.elements)
        if (startElement.overflow) {
          if (e.deltaY < startElement.offset) {
            window.scrollBy(0, e.deltaY)
          } else {
            this.isIntercept = true
            window.scrollBy(0, startElement.offset + ScrollYOffset)
          }
        } else {
          this.isIntercept = true
          this.elements[startElement.index + 1].scrollIntoView({ behavior: "smooth" })
        }
      } else if (e.deltaY < 0) {
        // up

        let endElement = indexEndAtWindow(this.elements)
        if (endElement.overflow) {
          if (e.deltaY < -endElement.offset) {
            window.scrollBy(0, e.deltaY)
          } else {
            this.isIntercept = true
            window.scrollBy(0, -endElement.offset)
          }
        } else {
          this.isIntercept = true
          this.elements[endElement.index - 1].scrollIntoView({ behavior: "smooth" })
        }
      }
    },
    scrollEnd(e) {
      console.log("scroll end");
      this.isIntercept = false
    },

    mousewhell(e) {
      e.preventDefault()
      window.clearTimeout(this.scrollTimer)
      this.scrollTimer = window.setTimeout(() => {
        this.scrollEnd(e)
      }, 300)
      if (300 < new Date() - this.scrollTime) {
        this.scrollStart(e)
      } else {
        this.scrollDo(e)
      }
      this.scrollTime = new Date()
    },
  },

  created() {
    document.addEventListener('scroll', this.scrollHandle, { passive: false })
    document.addEventListener('touchmove', this.scrollHandle, { passive: false });  //passive 参数不能省略，用来兼容ios和android
    document.addEventListener('wheel', this.mousewhell, { passive: false })
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
    document.removeEventListener('scroll', this.scrollHandle)
    document.removeEventListener('touchmove', this.scrollHandle)
    document.removeEventListener('wheel', this.mousewhell)
  },

  template: com_fullpage,
})