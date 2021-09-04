// vue_textmap.js

Vue.component('textmap', {

  props: ['scroller', 'content', 'range'],

  data: function () {
    return {
      documentScale: 0.111,
      documentHeight: 0,
      textmapHeight: 0,
      scrollbarHeight: 0,
      start: { x: 0, y: 0 },
      slider: {
        height: 56,
      },
      touch: false,
      top: 0,
      offset: 0,
    }
  },

  watch: {
    range(value, oldVal) {
      this.highlightRange(this.$el.firstChild, value.start, value.end)
    },

    scroller(value, oldVal) {
      oldVal.removeEventListener('scroll', this.scrollHandle)
      value.addEventListener('scroll', this.scrollHandle)
    }
  },

  methods: {
    init(e) {
      if (this.touch) return
      let documentHeight = this.getScrollHeight()
      this.documentHeight = documentHeight
      this.textmapHeight = this.$el.clientHeight
      this.scrollbarHeight = documentHeight * this.documentScale
      if (this.scrollbarHeight < this.textmapHeight) {
        this.textmapHeight = this.scrollbarHeight
      }
      this.slider.height = window.innerHeight * this.documentScale
    },

    onTouchTextmap(e) {
      if (this.touch) return

      let r = this.$el.getBoundingClientRect()
      this.init(e)

      let documentHeight = this.documentHeight
      let yAtTextmap = -this.offset + e.clientY - r.top - this.slider.height / 5.0
      let scrollY = yAtTextmap / this.scrollbarHeight * documentHeight
      this.scrollTo(scrollY)
    },

    startDrag(e) {
      this.init(e)
      this.touch = true
      this.start.x = e.clientX
      this.start.y = e.clientY
    },

    stopDrag(e) {
      if (this.touch) {
        this.start.x = e.clientX
        this.start.y = e.clientY
      }

      this.touch = false
    },

    doDrag(e) {
      if (!this.touch) return false

      let y = e.clientY
      let mapScale = this.textmapHeight / this.documentHeight
      let dict = y - this.start.y
      let scrollBy = dict / mapScale
      let scrollY = this.getScrollY()

      if (scrollY + scrollBy < 0) {
        this.scrollTo(0)
        return
      }

      this.scrollBy(scrollBy)

      this.start.x = e.clientX
      this.start.y = y
    },

    scrollHandle(e) {
      this.init(e)

      let y = this.getScrollY()
      let windowHeight = window.innerHeight
      let documentHeight = this.documentHeight
      let ratio = (documentHeight - windowHeight) / (this.textmapHeight - this.slider.height)
      let sliderTop = y / ratio
      this.top = sliderTop

      if (this.textmapHeight < this.scrollbarHeight) {
        ratio = (documentHeight - windowHeight) / (this.scrollbarHeight - this.textmapHeight)
        let offset = y / ratio
        this.offset = -offset
      }
    },

    getScrollHeight() {
      var temp = this.scroller
      if (temp === document) {
        return document.body.clientHeight
      }
      return temp.scrollHeight
    },

    getScrollY() {
      var temp = this.scroller
      if (temp === document) {
        return window.scrollY
      }
      return temp.scrollTop
    },

    scrollTo(y) {
      var temp = this.scroller
      if (temp === document) {
        temp = window
      }
      temp.scrollTo(0, y)
    },

    scrollBy(y) {
      var temp = this.scroller
      if (temp === document) {
        temp = window
      }
      temp.scrollBy(0, y)
    },

    selectTextRange(obj, start, stop) {
      var endNode, startNode = endNode = obj.firstChild

      startNode.nodeValue = startNode.nodeValue.trim();

      var range = document.createRange();
      range.setStart(startNode, start);
      range.setEnd(endNode, stop + 1);

      var sel = window.getSelection();
      sel.removeAllRanges();
      sel.addRange(range);
    },

    highlightRange(el, start, end) {
      if (!this.content) return

      if (end < start) {
        let temp = start
        start = end
        end = temp
      }
      var text = this.content
      var startText = text.substring(0, start)
      var highText = text.substring(start, end)
      var endText = text.substring(end)

      var startEl = document.createElement("span")
      startEl.innerText = startText
      var highEl = document.createElement("span")
      highEl.style.background = "var(--primary)"
      highEl.style.color = "var(--white)"
      highEl.innerText = highText
      var endEl = document.createElement("span")
      endEl.innerText = endText

      el.innerHTML = ''
      el.appendChild(startEl)
      el.appendChild(highEl)
      el.appendChild(endEl)
    },
  },

  created() {
    this.slider.height = window.innerHeight * this.documentScale

    document.addEventListener('mousemove', this.doDrag)
    document.addEventListener('mouseup', this.stopDrag)
    this.scroller.addEventListener('scroll', this.scrollHandle)
  },

  beforeDestroy() {
    document.removeEventListener('mousemove', this.doDrag)
    document.removeEventListener('mouseup', this.stopDrag)
    this.scroller.removeEventListener('scroll', this.scrollHandle)
  },

  template: com_textmap,
})