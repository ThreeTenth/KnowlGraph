// vue_textmap.js

Vue.component('textmap', {

  props: {
    content: String,
    cursor: {
      height: Number,
    },
  },

  data: function () {
    return {
      docHeight: 0,
      mapHeight: 0,
      scrollbarHeight: 0,
      start: { x: 0, y: 0 },
      touch: false,
      top: 0,
      offset: 0,
    }
  },

  computed: {
    // top: function () {
    //   console.log(this.scrollY);
    //   return 'top=' + this.scrollY + "px"
    // },
  },

  methods: {
    init(e) {
      if (this.touch) return
      this.docHeight = document.body.clientHeight
      this.mapHeight = this.$el.clientHeight
      this.scrollbarHeight = this.docHeight * 0.11
      if (this.scrollbarHeight < this.mapHeight) {
        this.mapHeight = this.scrollbarHeight
      }
    },

    scrollTo(e) {
      this.init(e)
      this.touch = true
      this.doDrag(e)
    },

    startDrag(e) {
      this.init(e)
      this.touch = true
      this.start.x = e.clientX
      this.start.y = e.clientY
    },

    stopDrag(e) {
      if (!this.touch) return false

      this.start.x = e.clientX
      this.start.y = e.clientY
      this.touch = false
    },

    doDrag(e) {
      if (!this.touch) return false

      let y = e.clientY
      let mapScale = this.mapHeight / this.docHeight
      let dict = y - this.start.y
      let scrollY = dict / mapScale

      window.scrollBy(0, scrollY)

      this.start.x = e.clientX
      this.start.y = y
    },

    scrollHandle(e) {
      this.init(e)

      let y = window.scrollY
      let mapScale = this.mapHeight / this.docHeight

      this.top = y * mapScale

      if (this.mapHeight < this.scrollbarHeight) {
        let barScale = this.scrollbarHeight / this.mapHeight
        this.offset = 0 - this.top * barScale / 2
      }
    },
  },

  created() {
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