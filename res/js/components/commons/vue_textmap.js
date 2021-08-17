// vue_textmap.js

Vue.component('textmap', {

  props: {
    content: String,
    range: {
      start: 0,
      end: 0,
    },
  },

  data: function () {
    return {
      docScale: 0.111,
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
  },

  methods: {
    init(e) {
      if (this.touch) return
      let documentHeight = document.body.clientHeight
      this.textmapHeight = this.$el.clientHeight
      this.scrollbarHeight = documentHeight * this.docScale
      if (this.scrollbarHeight < this.textmapHeight) {
        this.textmapHeight = this.scrollbarHeight
      }
      this.slider.height = window.innerHeight * this.docScale
    },

    scrollTo(e) {
      if (this.touch) return

      let r = this.$el.getBoundingClientRect()
      this.init(e)

      let documebtHeight = document.body.clientHeight
      let yAtTextmap = -this.offset + e.clientY - r.top - this.slider.height / 5.0
      let scrollY = yAtTextmap / this.scrollbarHeight * documebtHeight
      window.scrollTo(0, scrollY)
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
      let mapScale = this.textmapHeight / document.body.clientHeight
      let dict = y - this.start.y
      let scrollBy = dict / mapScale
      let scrollY = window.scrollY

      if (scrollY + scrollBy < 0) {
        window.scrollTo(0, 0)
        return
      }

      window.scrollBy(0, scrollBy)

      this.start.x = e.clientX
      this.start.y = y
    },

    scrollHandle(e) {
      this.init(e)

      let y = window.scrollY
      let windowHeight = window.innerHeight
      let documebtHeight = document.body.clientHeight
      let ratio = (documebtHeight - windowHeight) / (this.textmapHeight - this.slider.height)
      let sliderTop = y / ratio
      this.top = sliderTop

      if (this.textmapHeight < this.scrollbarHeight) {
        ratio = (documebtHeight - windowHeight) / (this.scrollbarHeight - this.textmapHeight)
        let offset = y / ratio
        this.offset = -offset
      }
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
      if (end < start) {
        let temp = start
        start = end
        end = temp
      }
      var text = this.content
      var startText = text.substring(0, start)
      var highText = '<span style="background: var(--primary); color: var(--white)">' +
        text.substring(start, end) +
        "</span>"
      var endText = text.substring(end)
      el.innerHTML = startText + highText + endText;
    },
  },

  created() {
    this.slider.height = window.innerHeight * this.docScale

    document.addEventListener('mousemove', this.doDrag)
    document.addEventListener('mouseup', this.stopDrag)
    document.addEventListener('scroll', this.scrollHandle)
  },

  beforeDestroy() {
    document.removeEventListener('mousemove', this.doDrag)
    document.removeEventListener('mouseup', this.stopDrag)
    document.removeEventListener('scroll', this.scrollHandle)
  },

  template: com_textmap,
})