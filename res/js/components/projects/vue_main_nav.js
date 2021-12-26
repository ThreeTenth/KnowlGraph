// vue_main_nav.js

Vue.component('main-nav', {

  props: {
    mobileHiddenMenu: Boolean,
    hiddenSearch: Boolean,
    openSearchShortcut: {
      type: Boolean,
      default: true,
    },
  },

  data: function () {
    return {
      profilePicture: getLink("icon"),
      showSearch: false,
    }
  },

  watch: {
    showSearch: function (val, oldVal) {
      if (val) {
        setTimeout(() => this.$refs.search.focus(), 0);
      }
    },
  },

  methods: {
    onSearch() {
      this.showSearch = true
    },
  },

  created() {
    document.onkeydown = (e) => {
      let key = window.event.keyCode
      if (key == 191 && !e.target.type && this.openSearchShortcut) {
        this.onSearch()
      }
    }
  },

  template: com_main_nav,
})