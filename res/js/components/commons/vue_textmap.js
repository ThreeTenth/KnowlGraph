// vue_textmap.js

Vue.component('textmap', {

  props: {
    content: String,
  },

  data: function () {
    return {
      docScale: 0.111,
      mapHeight: 0,
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

  computed: {
    // top: function () {
    //   console.log(this.scrollY);
    //   return 'top=' + this.scrollY + "px"
    // },
  },

  methods: {
    init(e) {
      if (this.touch) return
      let docHeight = document.body.clientHeight
      this.mapHeight = this.$el.clientHeight
      this.scrollbarHeight = docHeight * this.docScale
      if (this.scrollbarHeight < this.mapHeight) {
        this.mapHeight = this.scrollbarHeight
      }
      this.slider.height = window.innerHeight * this.docScale
    },

    scrollTo(e) {
      this.init(e)
      this.touch = true
      this.doDrag(e)
      this.touch = false
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
      let mapScale = this.mapHeight / document.body.clientHeight
      let dict = y - this.start.y
      let scrollY = dict / mapScale

      window.scrollBy(0, scrollY)

      this.start.x = e.clientX
      this.start.y = y
    },

    scrollHandle(e) {
      this.init(e)

      let y = window.scrollY
      let windowHeight = window.innerHeight
      let documebtHeight = document.body.clientHeight
      let ratio = (documebtHeight - windowHeight) / (this.mapHeight - this.slider.height)
      let sliderTop = y / ratio
      this.top = sliderTop

      if (this.mapHeight < this.scrollbarHeight) {
        ratio = (documebtHeight - windowHeight) / (this.scrollbarHeight - this.mapHeight)
        let offset = y / ratio
        this.offset = -offset
      }
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