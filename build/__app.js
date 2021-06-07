// languages.json, languages 
 const languages = [
  {"code":"sq","name":"Shqiptare", "direction": "ltr", "comment": "Albanian"},
  {"code":"ar","name":"عربى", "direction": "rtl", "comment": "Arabic"},
  {"code":"am","name":"አማርኛ", "direction": "ltr", "comment": "Amharic"},
  {"code":"az","name":"Azərbaycan dili", "direction": "ltr", "comment": "Azerbaijani"},
  {"code":"ga","name":"Gaeilge", "direction": "ltr", "comment": "Irish"},
  {"code":"et","name":"Eestlane", "direction": "ltr", "comment": "Estonian"},
  {"code":"or","name":"ଓଡିଆ", "direction": "ltr", "comment": "Oriya"},
  {"code":"eu","name":"Euskara", "direction": "ltr", "comment": "Basque"},
  {"code":"be","name":"Беларуская", "direction": "ltr", "comment": "Belarusian"},
  {"code":"bg","name":"български", "direction": "ltr", "comment": "Bulgarian"},
  {"code":"is","name":"Íslenska", "direction": "ltr", "comment": "Icelandic"},
  {"code":"pl","name":"Polskie", "direction": "ltr", "comment": "Polish"},
  {"code":"bs","name":"Bosanski", "direction": "ltr", "comment": "Bosnian"},
  {"code":"fa","name":"فارسی", "direction": "rtl", "comment": "Persian"},
  {"code":"tt","name":"Татар", "direction": "ltr", "comment": "Tatar"},
  {"code":"da","name":"dansk", "direction": "ltr", "comment": "Danish"},
  {"code":"de","name":"Deutsch", "direction": "ltr", "comment": "German"},
  {"code":"ru","name":"русский", "direction": "ltr", "comment": "Russian"},
  {"code":"fr","name":"français", "direction": "ltr", "comment": "French"},
  {"code":"tl","name":"Tagalog","direction": "ltr", "comment": "Tagalog"},
  {"code":"fi","name":"Suomalainen", "direction": "ltr", "comment": "Finnish"},
  {"code":"fy","name":"Westfrysk", "direction": "ltr", "comment": "Western Frisian"},
  {"code":"km","name":"ខ្មែរកណ្ដាល", "direction": "ltr", "comment": "Central Khmer"},
  {"code":"ka","name":"ქართველი", "direction": "ltr", "comment": "Georgian"},
  {"code":"gu","name":"ગુજરાતી", "direction": "ltr", "comment": "Gujarati"},
  {"code":"kk","name":"Қазақ", "direction": "ltr", "comment": "Kazakh"},
  {"code":"ht","name":"Kreyòl ayisyen", "direction": "ltr", "comment": "Haitian Creole"},
  {"code":"ko","name":"한국어", "direction": "ltr", "comment": "Korean"},
  {"code":"ha","name":"Hausa", "direction": "ltr", "comment": "Hausa"},
  {"code":"nl","name":"Nederlands", "direction": "ltr", "comment": "Dutch"},
  {"code":"ky","name":"Кыргызча", "direction": "ltr", "comment": "Kyrgyz"},
  {"code":"gl","name":"Galego", "direction": "ltr", "comment": "Galician"},
  {"code":"ca","name":"Català", "direction": "ltr", "comment": "Catalan"},
  {"code":"cs","name":"čeština", "direction": "ltr", "comment": "Czech"},
  {"code":"kn","name":"ಕನ್ನಡ", "direction": "ltr", "comment": "Kannada"},
  {"code":"hr","name":"Hrvatski", "direction": "ltr", "comment": "Croatian"},
  {"code":"ku","name":"Kurdî", "direction": "ltr", "comment": "Kurdish"},
  {"code":"la","name":"Latine", "direction": "ltr", "comment": "Latin"},
  {"code":"lv","name":"Latvietis", "direction": "ltr", "comment": "Latvian"},
  {"code":"lo","name":"ລາວ", "direction": "ltr", "comment": "Lao"},
  {"code":"lt","name":"Lietuvis", "direction": "ltr", "comment": "Lithuanian"},
  {"code":"lb","name":"Lëtzebuergesch", "direction": "ltr", "comment": "Luxembourgish"},
  {"code":"rw","name":"Ururimi rw'u Rwanda", "direction": "ltr", "comment": "Kinyarwanda"},
  {"code":"ro","name":"Română", "direction": "ltr", "comment": "Romanian"},
  {"code":"mg","name":"Teny malagasy", "direction": "ltr", "comment": "Malagasy"},
  {"code":"mt","name":"Malti", "direction": "ltr", "comment": "Maltese"},
  {"code":"mr","name":"मराठी", "direction": "ltr", "comment": "Marathi"},
  {"code":"ml","name":"മലയാളം", "direction": "ltr", "comment": "Malayalam"},
  {"code":"ms","name":"Bahasa Melayu", "direction": "ltr", "comment": "Malay"},
  {"code":"mk","name":"Македонски зборови", "direction": "ltr", "comment": "Macedonian"},
  {"code":"mi","name":"Reo Maori", "direction": "ltr", "comment": "Maori"},
  {"code":"mn","name":"Монгол", "direction": "ltr", "comment": "Mongolian"},
  {"code":"bn","name":"বাংলা", "direction": "ltr", "comment": "Bengali"},
  {"code":"my","name":"ဗမာ", "direction": "ltr", "comment": "Burmese"},
  {"code":"hmn","name":"Hmong", "direction": "ltr", "comment": "Hmong"},
  {"code":"xh","name":"isiXhosa", "direction": "ltr", "comment": "Xhosa"},
  {"code":"zu","name":"Ulimi lwesiZulu", "direction": "ltr", "comment": "Zulu"},
  {"code":"ne","name":"नेपाली", "direction": "ltr", "comment": "Nepali"},
  {"code":"no","name":"norsk", "direction": "ltr", "comment": "Norwegian"},
  {"code":"pa","name":"ਪੰਜਾਬੀ ਭਾਸ਼ਾ", "direction": "ltr", "comment": "Punjabi"},
  {"code":"pt","name":"Português", "direction": "ltr", "comment": "Portuguese"},
  {"code":"ps","name":"پښتو", "direction": "rtl", "comment": "Pashto"},
  {"code":"ny","name":"Chilankhulo cha Nyanja", "direction": "ltr", "comment": "Chichewa"},
  {"code":"ja","name":"日本語", "direction": "ltr", "comment": "Japanese"},
  {"code":"sv","name":"svenska", "direction": "ltr", "comment": "Swedish"},
  {"code":"sm","name":"Faasamoa", "direction": "ltr", "comment": "Samoan"},
  {"code":"sr","name":"Српски", "direction": "ltr", "comment": "Serbian"},
  {"code":"st","name":"Puo ya Sesotho sa Borwa", "direction": "ltr", "comment": "Sotho, Southern"},
  {"code":"si","name":"සිංහල", "direction": "ltr", "comment": "Sinhalese"},
  {"code":"eo","name":"Esperanto", "direction": "ltr", "comment": "Esperanto"},
  {"code":"sk","name":"Slovák", "direction": "ltr", "comment": "Slovak"},
  {"code":"sl","name":"Slovenščina", "direction": "ltr", "comment": "Slovenian"},
  {"code":"sw","name":"Kiswahili", "direction": "ltr", "comment": "Swahili"},
  {"code":"gd","name":"Gàidhlig na h-Alba", "direction": "ltr", "comment": "Scottish Gaelic"},
  {"code":"ceb","name":"Sugbuanon", "direction": "ltr", "comment": "Cebuano"},
  {"code":"so","name":"Soomaali", "direction": "ltr", "comment": "Somali"},
  {"code":"tg","name":"Тоҷикӣ", "direction": "ltr", "comment": "Tajik"},
  {"code":"te","name":"తెలుగు భాష", "direction": "ltr", "comment": "Telugu"},
  {"code":"ta","name":"தமிழ்", "direction": "ltr", "comment": "Tamil"},
  {"code":"th","name":"ไทย", "direction": "ltr", "comment": "Thai"},
  {"code":"tr","name":"Türk", "direction": "ltr", "comment": "Turkish"},
  {"code":"tk","name":"Türkmenler", "direction": "ltr", "comment": "Turkmen"},
  {"code":"cy","name":"Cymraeg", "direction": "ltr", "comment": "Welsh"},
  {"code":"ug","name":"ئۇيغۇر", "direction": "rtl", "comment": "Uighur"},
  {"code":"ur","name":"اردو", "direction": "rtl", "comment": "Urdu"},
  {"code":"uk","name":"Український", "direction": "ltr", "comment": "Ukrainian"},
  {"code":"uz","name":"O'zbek", "direction": "ltr", "comment": "Uzbek"},
  {"code":"es","name":"Español", "direction": "ltr", "comment": "Spanish"},
  {"code":"iw","name":"עִברִית", "direction": "rtl", "comment": "Hebrew"},
  {"code":"el","name":"Νέα Ελληνικά", "direction": "ltr", "comment": "Greek, Modern (1453-)"},
  {"code":"haw","name":"Ōlelo Hawaiʻi", "direction": "ltr", "comment": "Hawaiian"},
  {"code":"sd","name":"سنڌي ٻولي", "direction": "rtl", "comment": "Sindhi"},
  {"code":"hu","name":"Magyar", "direction": "ltr", "comment": "Hungarian"},
  {"code":"sn","name":"Mutauro weSchona", "direction": "ltr", "comment": "Shona"},
  {"code":"hy","name":"հայերեն", "direction": "ltr", "comment": "Armenian"},
  {"code":"ig","name":"Asusu Ibo", "direction": "ltr", "comment": "Igbo"},
  {"code":"it","name":"italiano", "direction": "ltr", "comment": "Italian"},
  {"code":"yi","name":"יידיש", "direction": "rtl", "comment": "Yiddish"},
  {"code":"hi","name":"हिंदी", "direction": "ltr", "comment": "Hindi"},
  {"code":"su","name":"Sundanis", "direction": "ltr", "comment": "Sundanese"},
  {"code":"id","name":"bahasa Indonesia", "direction": "ltr", "comment": "Indonesian"},
  {"code":"jw","name":"Wong jawa", "direction": "ltr", "comment": "Javanese"},
  {"code":"en","name":"English", "direction": "ltr", "comment": "English"},
  {"code":"yo","name":"Ede Yoruba", "direction": "ltr", "comment": "Yoruba"},
  {"code":"vi","name":"Tiếng Việt", "direction": "ltr", "comment": "Vietnamese"},
  {"code":"zh","name":"中文", "direction": "ltr", "comment": "Chinese"}
]

// vue_autocomplet.js

Vue.component('autocomplet', {
  props: ["placeholder", "source"],

  data: function () {
    return {
      input: "",
      items: [],
      auto: false,
      selected: 0,
      pos: 'down',
      watchTime: 0,
    }
  },

  watch: {
    items: function (val, oldVal) {
      let len = val.length

      if (0 == len) return

      this.$refs.selection.style.maxHeight = ''
      this.pos = "down"

      this.$nextTick(() => {
        let selection = this.$refs.selection
        let rect = selection.getBoundingClientRect()
        let windowHeight = document.documentElement.clientHeight
        let top = rect.top
        let preHeight = rect.height

        if (top + preHeight < windowHeight) {
          this.pos = "down"
        } else {
          this.pos = "up"

          selection.style.maxHeight = rect.top + 'px'
        }
      })
    },
  },

  mounted() {
    // this.__update()
  },

  activated() {
    // this.__update()
  },

  methods: {
    onSelect(item) {
      this.$emit('item', item, val => {
        this.input = val
      })
      this.$emit('select', item)
      this.items = []
      this.auto = false
    },

    onKeyTab() {
      console.log("on tab")
    },

    onHotKey(e) {
      switch (e.keyCode) {
        case 38:  // up
          if (0 == this.items.length) return
          if (0 < this.selected) {
            this.selected -= 1
          } else {
            this.selected = this.items.length - 1
          }
          this.$refs.selection.children[this.selected].scrollIntoView()
          e.preventDefault()
          break;
        case 9:   // tab
        case 40:  // down
          if (0 == this.items.length) return
          if (this.items.length <= this.selected + 1) {
            this.selected = 0
          } else {
            this.selected += 1
          }
          this.$refs.selection.children[this.selected].scrollIntoView()
          e.preventDefault()
          break;
        case 13:  // enter
          if (0 == this.items.length) return
          this.onSelect(this.items[this.selected])
          e.preventDefault()
          break;
      }
    },

    onChanged(e) {
      // console.log("onChanged", e.code, e.key, e.keyCode, e.target.value)

      switch (e.keyCode) {
        case 38:  // up
        case 9:   // tab
        case 40:  // down
        case 13:  // enter
          return
      }

      const items = []
      const regex = new RegExp('^' + e.target.value)
      this.source.forEach(item => {
        let value = item
        this.$emit('item', item, val => {
          value = val
        })

        if (regex.test(value)) {
          items.push(item)
        }
      });

      this.items = items
      this.selected = 0
      this.auto = true
      // this.__update()
    },

    // __update() {
    //   let len = this.items.length

    //   if (0 == len) return

    //   let autocomplete = this.$refs.autocomplete
    //   let bottom = autocomplete.getBoundingClientRect().bottom

    //   let preHeight = len * 30

    //   if (preHeight < bottom) {
    //     this.pos = "down"
    //   } else {
    //     this.pos = "up"
    //   }
    // },
  },

  template: com_autocomplet,
})

// vue_dropdown.js

Vue.component('dropdown', {
  props: ['ignore'],

  data: function () {
    return {
      seen: false,
      pos: 'down',
    }
  },
  computed: {
  },
  mounted() {
    this.__update()
  },
  activated() {
    this.__update()
  },
  methods: {
    toggle() {
      if (this.seen) {
        return this.__hide()
      }
      return this.__show()
    },

    __show() {
      this.seen = true
      setTimeout(() => document.addEventListener('click', this.__hide), 0);
    },

    __hide(e) {
      if (!e) return
      if (this.ignore) {
        var path = e.path
        var hide = true
        for (let index = 0; index < path.length; index++) {
          const element = path[index];
          if (element.className == "drop-content") {
            hide = false
            break
          }
        }
        if (hide) {
          this.seen = false
          document.removeEventListener('click', this.__hide);
        }
      } else {
        this.seen = false
        document.removeEventListener('click', this.__hide);
      }
    },

    __update() {
      let top = this.$refs.dropdown.getBoundingClientRect().top
      if (top < 240) {
        this.pos = 'down'
      } else {
        this.pos = 'up'
      }
    }
  },

  template: com_dropdown,
})

// vue_modal.js

Vue.component('modal', {
  model: {
    prop: 'seen',
    event: 'toggle'
  },

  props: {
    seen: Boolean
  },

  watch: {
    seen: function (val, oldVal) {
      if (val) {
        this.$nextTick(() => {
          const dialog = this.$refs['dialog-content']
          const target = document.body
          target.insertBefore(dialog, target.lastChild)
        })
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = ''
      }
    },
  },

  beforeDestroy() {
    if (this.$el && this.$el.parentNode) {
      this.$el.parentNode.removeChild(this.$el)
    }
  },

  methods: {
    close() {
      this.$emit('toggle', false)
    },
  },

  // template: com_modal,

  render(h) {
    const content = h('div', {
      staticClass: 'vs-dialog__content',
      // class: {
      //   notFooter: !this.$slots.footer
      // }
    }, [
      this.$slots.default
    ])

    const dialog = h('div', {
      staticClass: 'modal-content',
      // style: {
      //   width: this.width
      // },
      // class: {
      //   'vs-dialog--fullScreen': this.fullScreen,
      //   'vs-dialog--rebound': this.rebound,
      //   'vs-dialog--notPadding': this.notPadding,
      //   'vs-dialog--square': this.square,
      //   'vs-dialog--autoWidth': this.autoWidth,
      //   'vs-dialog--scroll': this.scroll,
      //   'vs-dialog--loading': this.loading,
      //   'vs-dialog--notCenter': this.notCenter,
      // }
    }, [
      // this.loading && loading,
      // !this.notClose && close,
      // this.$slots.header && header,
      this.$slots.default,
      // this.$slots.footer && footer
    ])

    const dialogContent = h('div', {
      staticClass: 'modal',
      ref: 'dialog-content',
      // class: {
      //   blur: this.blur,
      //   fullScreen: this.fullScreen,
      // },
      on: {
        click: (evt) => {
          if (!evt.target.closest('.modal-content') && !this.preventClose) {
            this.$emit('toggle', !this.seen)
            this.$emit('close')
          }
        }
      }
    }, [
      dialog
    ])

    const modalTransition = h('transition', {
      props: {
        name: 'vs-dialog'
      },
    }, [this.seen && dialogContent])

    return h('span', [modalTransition])
  },
})

// vue_time.js

Vue.component('the-time', {
  props: ['datetime'],

  computed: {
    fulltime: function () {
      let lang = this.i18n.__code
      var options = { weekday: "long", year: "numeric", month: "short", day: "numeric", hour: "numeric", minute: "numeric", hour12: false };
      return new Date(this.datetime).toLocaleString(lang, options)
    },
    timetonow: function () {
      let now = new Date()
      let timer = (now - new Date(this.datetime)) / 1000
      let tip = ''

      if (timer <= 0) {
        tip = this.i18n.JustNow
      } else if (Math.floor(timer / 60) <= 0) {
        tip = this.i18n.JustNow
      } else if (Math.floor(timer < 3600)) {
        tip = sprintf(this.i18n.MinutesAgo, Math.floor(timer / 60))
      } else if (timer >= 3600 && timer < 86400) {
        tip = sprintf(this.i18n.HoursAgo, Math.floor(timer / 3600))
      } else if (timer / 86400 <= 31) {
        tip = sprintf(this.i18n.DaysAgo, Math.floor(timer / 86400))
      } else {
        tip = this.fulltime
      }
      return tip
    },
  },

  template: com_the_time,
})

// vue word input

const WordStateIdle = 0
const WordStateMayDel = 1
const WordStatePreDel = 2
const WordStateSafe = 3

Vue.component('words-input', {
  props: ['words', 'keywords'],

  data: function () {
    return {
      auto: false,
      values: this.keywords ? this.keywords : [],
      input: "",
      newInput: "",
      items: [],
      deleteMark: null,
      __isDelProcess: false,
    }
  },

  methods: {
    onKeyup: function (event) {
      if ("Process" === event.key) {
        // 针对中文输入法的特殊处理
        // 在进入中文输入时，
        // vue 无法观测到正在输入的内容变化，
        // 但可以监听到按键变化。
        // 所以如果当前内容为空，
        // 但用户正在使用中文输入法输入内容，
        // 如果使用了删除键，
        // 这个删除事件就会被捕获，
        // 那么就会出现错误的删除操作。
        // 在 Vue 中，如果用户正在进行中文输入，
        // 那么每次按键事件都被前后触发两次，
        // 第一次是 key 始终为 "Process" 的按键事件，
        // 第二次是正常的按键事件，
        // 通过这个，可以过滤不需要的按键。
        // 在这里，就是删除键。

        this.$data.__isDelProcess = true
      }
    },
    onRemoveWord: function (word) {
      const index = this.values.indexOf(word)
      if (index > -1) {
        this.values.splice(index, 1)
      }
    },
    onAdd: function () {
      var word = this.input.replace(/[\,\， ]*$/g, '')
      if ("" == word) {
        this.input = ""
        // todo Waring toast
        return
      }

      var words = this.values
      for (let index = 0; index < words.length; index++) {
        const element = words[index];
        if (element == word) {
          this.input = ""
          // todo Waring toast
          return
        }
      }

      this.values.push(word)
      this.input = ""
    },
    onDel: function () {
      if (this.$data.__isDelProcess) {
        this.$data.__isDelProcess = false
        return
      }
      if (0 < this.newInput.length) return
      if (null == this.deleteMark) {
        cancelDeleteWordTimeoutID = window.setTimeout(() => { this.deleteMark = null }, 2000)
        this.deleteMark = this.values.length - 1
      } else {
        this.values.splice(this.deleteMark, 1)
        window.clearTimeout(cancelDeleteWordTimeoutID)
        this.deleteMark = null
      }
    },
  },

  watch: {
    input: function (val, oldVal) {
      window.setTimeout(() => { this.newInput = val }, 100)
      if ("" == val) {
        this.auto = false
      } else {
        this.deleteMark = null
        if (this.words) {
          if (2 <= val.length) {
            const items = []
            const regex = new RegExp('^' + val)
            this.words.forEach(item => {
              if (regex.test(item)) {
                items.push(item)
              }
            });
            this.items = items
            this.auto = 0 != items.length
          } else {
            this.auto = false
          }
        }
      }
    },
  },

  template: com_words_input,
})

// vue_new_node.js

Vue.component('new-node', {
  data: function () {
    return {
      selected: null,
      value: "",
    }
  },

  methods: {
    onSelectPrev: function (archive) {
      this.selected = this.selected == archive ? null : archive
    },

    onSelectWord(word) {
      this.value = word
    },

    getWordName(word, callback) {
      callback(word.name)
    }
  },

  template: com_new_node,
})

// archive.js

const Archive = {
  data: function () {
    return {
      status: "achiveSelf",
    }
  },

  created() {
    this.status = this.$router.history.current.name
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.status = this.$router.history.current.name
  },

  template: fgm_archive,
}

function getArchiveArticles(status) {
  return {
    data: function () {
      return {
        __original: [],
      }
    },

    computed: {
      articles: function () {
        let original = this.$data.__original
        if (original == undefined) return []
        var _articles = []
        for (let index = 0; index < original.length; index++) {
          const version = original[index];
          const article = version.edges.article
          let title = version.title ? version.title : ""
          let code = title
          if (!code) {
            code = version.gist
          }
          code = encodeURLTitle(code)

          _articles[index] = {
            id: article.id,
            status: article.status,
            title: title,
            gist: version.gist,
            code: code,
            created_at: version.created_at
          }
        }
        return _articles
      }
    },

    created() {
      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/assets", { status: status }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    methods: {
      onArticle(i) {
        let id = this.articles[i].id
        let code = this.articles[i].code
        router.push({
          name: 'article', params: {
            id: id,
            code: code
          }
        })
      }
    },

    template: fgm_archive_articles,
  }
}

const archive_router = [
  { path: '', name: 'achiveSelf', component: getArchiveArticles("self"), props: true },
  { path: 'star', name: 'achiveStar', component: getArchiveArticles("star"), props: true },
  { path: 'watch', name: 'achiveWatch', component: getArchiveArticles("watch"), props: true },
  { path: 'readlist', name: 'achiveReadlist', component: getArchiveArticles("browse"), props: true },
]


// article.js

const Article = {
  props: ['id', 'code'],

  data: function () {
    return {
      emoji: {
        "up": "👍", "down": "👎", "laugh": "😄", "hooray": "🎉", "confused": "😕", "heart": "❤️", "rocket": "🚀", "eyes": "👀",
      },
      reactions: [
        { "name": "Up", "value": "up", "emoji": "👍" },
        { "name": "Down", "value": "down", "emoji": "👎" },
        { "name": "Laugh", "value": "laugh", "emoji": "😄" },
        { "name": "Hooray", "value": "hooray", "emoji": "🎉" },
        { "name": "Confused", "value": "confused", "emoji": "😕" },
        { "name": "Heart", "value": "heart", "emoji": "❤️" },
        { "name": "Rocket", "value": "rocket", "emoji": "🚀" },
        { "name": "Eyes", "value": "eyes", "emoji": "👀" },
      ],
      __original: {},
      isNewNode: false,
    }
  },

  computed: {
    article: function () {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]
      let content = version.edges.content

      var title, gist, code = seo(version.title, version.gist)

      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      var body = converter.makeHtml(content.body);
      var keywords = version.edges.keywords

      var lang = userLang
      languages.forEach(element => {
        if (element.code == version.lang) {
          lang = element
        }
      });

      var reactions = article.edges.reactions
      var assets = article.edges.assets
      var nodes = article.edges.nodes
      var star = false
      var watch = false
      var private = false
      var browseCount = 0
      if (assets) {
        for (let index = 0; index < assets.length; index++) {
          const element = assets[index];
          if (element.status == "star") {
            star = true
          } else if (element.status == "watch") {
            watch = true
          } else if (element.status == "self") {
            private = true
          } else if (element.status == "browse") {
            browseCount++
          }
        }
      }

      // sort by value
      reactions.sort(function (a, b) {
        return (b.count - a.count)
      });

      return {
        id: article.id,
        status: article.status,
        title: title,
        gist: gist,
        body: body,
        lang: lang,
        reactions: reactions,
        keywords: keywords,
        created_at: version.created_at,
        code: code,
        star: star,
        watch: watch,
        private: private,
        browseCount: browseCount,
        nodes: nodes,
      }
    }
  },

  methods: {
    onEditArticle() {
      let article = this.$data.__original
      if (undefined == article.edges) return

      let version = article.edges.versions[0]

      axios({
        method: "PUT",
        url: queryRestful("/v1/article/edit", { id: version.id }),
      }).then(function (resp) {
        router.push({ name: 'editDraft', params: { id: resp.data.id, __draft: resp.data } })
      }).catch(function (resp) {
        console.log(resp.status, resp.data)
      })
    },

    onPickReaction(reaction) {
      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/reaction", { articleId: this.article.id, reaction: reaction }),
      }).then(function (resp) {
        let reactions = _this.article.reactions
        let ok = false
        reactions = reactions == undefined ? [] : reactions
        for (let index = 0; index < reactions.length; index++) {
          const reac = reactions[index];
          if (reac.status == reaction) {
            reac.count += 1
            ok = true
            break
          }
        }

        if (!ok) {
          reactions.push(resp.data)
        }

        _this.$data.__original.edges.reactions = reactions
      }).catch(function (resp) {
        console.error(resp.status, resp.data)
      })
    },

    onNewNode() {
      this.isNewNode = !this.isNewNode
    },

    onStar() {
      if (this.user.logined) {
        this.__putAsset('star')
      } else {
        router.push({ name: 'login' })
      }
    },

    onWatch() {
      if (this.user.logined) {
        this.__putAsset('watch')
      } else {
        router.push({ name: 'login' })
      }
    },

    switchLang() { },

    __putAsset(status) {
      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/asset", { articleId: this.article.id, status: status }),
      }).then(function (resp) {

      }).catch(function (resp) {
        if (status == "browse") return
        _this.$vs.notification({
          color: 'danger',
          position: 'bottom-center',
          title: "Failure",
          text: resp.data
        })
      })
    },

    __load(id) {
      if (this.$data.__original.id == this.id) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/article", { id: id }),
      }).then(function (resp) {
        _this.$data.__original = resp.data
        if (!_this.code) {
          router.replace({
            name: 'article', params: {
              id: _this.id,
              code: _this.article.code,
            }
          })
        }
        if (_this.user.logined) {
          _this.$nextTick(() => {
            _this.__putAsset("browse")
          })
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.__load(to.params.id)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id)
  },

  template: fgm_article,
}

// draft_history.js

const DraftHistories = {
  props: ['id'],

  data: function () {
    return {
      snapshots: [],
    }
  },

  methods: {
    onHistory: function (snapshot) {
      router.push({ name: 'draftHistory', params: { id: this.id, hid: snapshot.id, __snapshot: snapshot } })
    },

    __load(id) {
      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: id, needHistory: true }),
      }).then(function (resp) {
        _this.snapshots = resp.data.edges.snapshots
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    if (from.name === "draftHistory") {
      next()
    } else {
      next(vm => {
        vm.__load(to.params.id)
      })
    }
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id)
  },

  template: fgm_draft_histories,
}

const DraftHistory = {
  props: ['id', 'hid', '__snapshot'],

  data: function () {
    return {
      snapshot: this.__snapshot,
    }
  },

  methods: {
    onHistory: function () {
    },
    onBack: function () {
      if (3 <= window.history.length) {
        router.go(-1)
      } else {
        router.push({ name: 'draftHistories', params: { id: this.id } })
      }
    },
    __load(id, hid) {
      if (null != this.snapshot) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: id, historyID: hid }),
      }).then(function (resp) {
        let snapshots = resp.data.edges.snapshots
        if (snapshots) {
          _this.snapshot = snapshots[0]
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next(vm => {
      vm.snapshot = to.params.__snapshot
      vm.__load(to.params.id, to.params.hid)
    })
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.snapshot = to.params.__snapshot
    this.__load(to.params.id, to.params.hid)
  },

  template: fgm_draft_history,
}

// drafts.js

const Drafts = {
  data: function () {
    return {
      __original: [],
    }
  },

  computed: {
    drafts: function () {
      let original = this.$data.__original
      var _drafts = []
      for (let index = 0; index < original.length; index++) {
        const element = original[index];
        let snapshot = element.edges.snapshots[0]
        let body = removeMarkdown(snapshot.body)
        if (200 < body.length) {
          body = body.substr(0, 120).trim() + '...'
        }

        _drafts[index] = { body: body, created_at: snapshot.created_at }
      }
      return _drafts
    }
  },

  created() {
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/drafts"),
    }).then(function (resp) {
      _this.$data.__original = resp.data
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  methods: {
    onDraft(i) {
      let draft = this.$data.__original[i]
      router.push({ name: 'editDraft', params: { id: draft.id, __draft: draft } })
    }
  },

  template: fgm_user_drafts,
}

// index.js

const Index = {
  data: function () {
    return {
      __original: [],
    }
  },

  computed: {
    articles: function () {
      let original = this.$data.__original
      if (!original) return []
      var _articles = []
      for (let index = 0; index < original.length; index++) {
        const version = original[index];
        const article = version.edges.article
        let title = version.title ? version.title : ""
        let code = title
        if (!code) {
          code = version.gist
        }
        code = encodeURLTitle(code)

        _articles[index] = {
          id: article.id,
          status: article.status,
          title: title,
          gist: version.gist,
          code: code,
          created_at: version.created_at
        }
      }
      return _articles
    }
  },

  created() {
    let _this = this
    axios({
      method: "GET",
      url: queryRestful("/v1/articles", { limit: 10, offset: 0 }),
    }).then(function (resp) {
      _this.$data.__original = resp.data
    }).catch(function (resp) {
      console.log(resp)
    })
  },

  methods: {
    onArticle(i) {
      let id = this.articles[i].id
      let code = this.articles[i].code
      router.push({
        name: 'article', params: {
          id: id,
          code: code
        }
      })
    }
  },

  template: fgm_home,
}

// new_article.js

var postChangedTimeoutID;

const NewArticle = {
  beforeRouteEnter(to, from, next) {
    axios({
      method: "PUT",
      url: queryRestful("/v1/article", { status: to.query.status }),
    }).then(function (resp) {
      router.push({ name: 'editDraft', params: { id: resp.data.id, __draft: resp.data } })
    }).catch(function (resp) {
      console.log(resp.status, resp.data)
    })
  },

  template: "",
}

const EditDraft = {
  props: ['id', "__draft"],

  data: function () {
    return {
      draft: this.__draft,
      __last: { body: "" },
    }
  },

  computed: {
    body: {
      get: function () {
        if (this.draft && this.draft.edges.snapshots)
          return this.draft.edges.snapshots[0].body
        return ""
      },
      set: function (newValue) {
        this.draft.edges.snapshots[0].body = newValue
      }
    },

    status: function () {
      if (!this.draft) { return "" }
      return this.draft.edges.article.status
    },
  },

  methods: {
    onHistories() {
      router.push({ name: 'draftHistories', params: { id: this.id } })
    },

    onPublish() {
      router.push({ name: 'publishArticle', params: { id: this.id, __draft: this.draft } })
    },

    onChanged: function () {
      let _this = this
      window.clearTimeout(postChangedTimeoutID)
      postChangedTimeoutID = window.setTimeout(function () { _this.onSaveDraft() }, 2000)
    },

    onBlur: function () {
      window.clearTimeout(postChangedTimeoutID)
      this.onSaveDraft()
    },

    onSaveDraft: function () {
      // console.trace()
      let content = this.draft.edges.snapshots[0]
      if ("" === content.body) return
      if (this.__last && content.body && content.body === this.__last) {
        return
      }

      let _this = this
      axios({
        method: "PUT",
        url: queryRestful("/v1/article/content"),
        data: {
          body: content.body,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        _this.draft.edges.snapshots[0] = resp.data
        _this.__setLast(_this.draft)
        // todo the new article have content and notify drafts page
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __load(__id, __draft) {
      if (__draft) {
        this.__setDraft(__draft)
        this.__setLast(__draft)
        return
      }

      if (this.draft && __id == this.draft.id) {
        this.__setLast(this.draft)
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: __id }),
      }).then(function (resp) {
        _this.__setDraft(resp.data)
        _this.__setLast(resp.data)
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __setDraft(__draft) {
      this.draft = __draft
      if (!this.draft.edges.snapshots) {
        this.draft.edges.snapshots = [{ body: "" }]
      }
    },

    __setLast(__draft) {
      this.__last = __draft.edges.snapshots[0].body
    }
  },

  created() {
    this.__load(this.id, this.__draft)
  },

  activated() {
    this.__load(this.id, this.__draft)
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id, to.params.__draft)
  },

  template: fgm_new_article,
}

// personalize.js

const Personalize = {
  computed: {
    currentLang: function () {
      return sprintf(this.i18n.CurrentUserLang, this.user.lang.name)
    },
  },

  methods: {
    onSelectUserLang: function (lang) {
      // const lang = event.target.value
      const old = this.user.lang
      const langCode = lang.code

      this.user.lang = lang
      setUserLang(langCode)

      if (langCode == defaultLang.__code) {
        this.__setUserLangSuccessed(defaultLang)
        return
      } else if (i18ns.has(langCode)) {
        this.__setUserLangSuccessed(i18ns.get(langCode))
        return
      } else {
        let strings = getI18nStrings(langCode)
        if (null != strings && strings.__version == defaultLang.__version) {
          i18ns.set(langCode, strings)
          this.__setUserLangSuccessed(strings)
          return
        }
        removeI18nStrings(langCode)
      }

      let _this = this

      axios({
        method: "GET",
        url: queryStatic("/static/strings/strings-" + langCode + ".json"),
      }).then(function (resp) {
        i18ns.set(langCode, resp.data)
        setI18nStrings(langCode, resp.data)
        _this.__setUserLangSuccessed(resp.data)
      }).catch(function (resp) {
        _this.user.lang = old
        setUserLang(old.code)
        _this.$vs.notification({
          color: 'danger',
          position: 'top-right',
          title: _this.i18n.SetUserLangFailure,
        })
      })
    },

    __setUserLangSuccessed(stringI18N) {
      Object.assign(i18n, stringI18N)
      this.$vs.notification({
        color: 'success',
        position: 'top-right',
        title: this.i18n.SetUserLangSuccess,
      })
    }
  },

  template: fgm_personalize
}

// publish_article.js

const PublishArticle = {
  props: ['id', '__draft'],

  data: function () {
    return {
      draft: {},
      cover: "",
      title: "",
      gist: "",
      versionName: "",
      comment: "",
      lang: getUserLang(),
      keywords: [],
      words: [],
      languages: languages,
      status: "public",
    }
  },

  watch: {
    draft: function (val, oldVal) {
      if (!val) return ""

      let content = val.edges.snapshots[0].body
      let text = removeMarkdown(content)
      let gist

      if (200 <= text.length) {
        gist = text.substr(0, 120).trim() + '...'
      } else {
        gist = text
      }
      
      gist = gist.replace(/[\r|\n]/g, ' ')

      this.status = val.edges.article.status

      let found = content.match(/^# (.*)/m)

      if (found) {
        this.title = found[1]
        this.gist = gist.replace(found[1], '').trim()
        return
      }

      found = content.match(/^#+ (.*)/m)

      if (found) {
        this.title = found[1]
        this.gist = gist.replace(found[1], '').trim()
        return
      }
    },
  },

  methods: {
    onPublish() {
      var _this = this

      axios({
        method: "PUT",
        url: queryRestful("/v1/publish/article"),
        data: {
          name: this.versionName,
          comment: this.comment,
          cover: this.cover,
          title: this.title,
          gist: this.gist,
          lang: this.lang,
          keywords: this.keywords,
          draft_id: this.draft.id,
        },
      }).then(function (resp) {
        // todo this content is published and notify home page and drafts page
        if (204 == resp.status) {
          router.push({ path: '/' })
          return
        }
        router.push({
          name: 'article', params: {
            id: resp.data.edges.article.id,
            code: encodeURLTitle(_this.title)
          }
        })
      }).catch(function (resp) {
        console.log(resp)
      })
    },

    __load(__id, __draft) {
      if (__draft) {
        this.draft = __draft
        return
      }

      if (this.draft && __id == this.draft.id) {
        return
      }

      let _this = this
      axios({
        method: "GET",
        url: queryRestful("/v1/draft", { id: __id }),
      }).then(function (resp) {
        _this.draft = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })

      axios({
        method: "GET",
        url: queryRestful("/v1/words"),
      }).then(function (resp) {
        _this.words = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  },

  created() {
    this.__load(this.id, this.__draft)
  },

  activated() {
    this.__load(this.id, this.__draft)
  },

  beforeRouteEnter(to, from, next) {
    if (!logined) {
      router.push({ name: "login" })
      return
    }

    next()
  },

  beforeRouteUpdate(to, from, next) {
    next()
    this.__load(to.params.id, to.params.__draft)
  },

  template: fgm_publish_article,
}

// rest.js

function queryRestful(url, data = undefined) {
  url = restfulDomain + url
  return encodeQueryData(url, data)
}

function queryStatic(url, data = undefined) {
  url = staticDomain + url
  return encodeQueryData(url, data)
}

function encodeQueryData(url, data = undefined) {
  if (undefined == data) return url
  const ret = [];
  for (let d in data)
    ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));

  if (0 == ret.length) return url

  return url + "?" + ret.join('&');
}

// utils.js

function seo(title, gist) {
  let vt = title
  if (!vt) {
    vt = gist
    if (60 < gist.length) {
      let idx = vt.search(/[.?!。？！]/g)
      if (1 < idx) {
        vt = gist.substr(0, idx)
      }
    }
  }

  return vt, gist, encodeURLTitle(vt)
}

function encodeURLTitle(title) {
  return title.trim()
    .replace(/[;,/?:@&=+$_.!~*'()# \n]+/g, '-')
    .replace(/-$/g, '')
    .toLowerCase()
}

// by @Stian Grytøyr's https://github.com/stiang/remove-markdown
function removeMarkdown(md, options) {
  options = options || {};
  options.listUnicodeChar = options.hasOwnProperty('listUnicodeChar') ? options.listUnicodeChar : false;
  options.stripListLeaders = options.hasOwnProperty('stripListLeaders') ? options.stripListLeaders : true;
  options.gfm = options.hasOwnProperty('gfm') ? options.gfm : true;
  options.useImgAltText = options.hasOwnProperty('useImgAltText') ? options.useImgAltText : true;

  var output = md || '';

  // Remove horizontal rules (stripListHeaders conflict with this rule, which is why it has been moved to the top)
  output = output.replace(/^(-\s*?|\*\s*?|_\s*?){3,}\s*$/gm, '');

  try {
    if (options.stripListLeaders) {
      if (options.listUnicodeChar)
        output = output.replace(/^([\s\t]*)([\*\-\+]|\d+\.)\s+/gm, options.listUnicodeChar + ' $1');
      else
        output = output.replace(/^([\s\t]*)([\*\-\+]|\d+\.)\s+/gm, '$1');
    }
    if (options.gfm) {
      output = output
        // Header
        .replace(/\n={2,}/g, '\n')
        // Fenced codeblocks
        .replace(/~{3}.*\n/g, '')
        // Strikethrough
        .replace(/~~/g, '')
        // Fenced codeblocks
        .replace(/`{3}.*\n/g, '');
    }
    output = output
      // Remove HTML tags
      .replace(/<[^>]*>/g, '')
      // Remove setext-style headers
      .replace(/^[=\-]{2,}\s*$/g, '')
      // Remove footnotes?
      .replace(/\[\^.+?\](\: .*?$)?/g, '')
      .replace(/\s{0,2}\[.*?\]: .*?$/g, '')
      // Remove images
      .replace(/\!\[(.*?)\][\[\(].*?[\]\)]/g, options.useImgAltText ? '$1' : '')
      // Remove inline links
      .replace(/\[(.*?)\][\[\(].*?[\]\)]/g, '$1')
      // Remove blockquotes
      .replace(/^\s{0,3}>\s?/g, '')
      // Remove reference-style links?
      .replace(/^\s{1,2}\[(.*?)\]: (\S+)( ".*?")?\s*$/g, '')
      // Remove atx-style headers
      .replace(/^(\n)?\s{0,}#{1,6}\s+| {0,}(\n)?\s{0,}#{0,} {0,}(\n)?\s{0,}$/gm, '$1$2$3')
      // Remove emphasis (repeat the line to remove double emphasis)
      .replace(/([\*_]{1,3})(\S.*?\S{0,1})\1/g, '$2')
      .replace(/([\*_]{1,3})(\S.*?\S{0,1})\1/g, '$2')
      // Remove code blocks
      .replace(/(`{3,})(.*?)\1/gm, '$2')
      // Remove inline code
      .replace(/`(.+?)`/g, '$1')
      // Replace two or more newlines with exactly two? Not entirely sure this belongs here...
      .replace(/\n{2,}/g, '\n\n');
  } catch (e) {
    console.error(e);
    return md;
  }
  return output;
};

function getUserLang() {
  var language = Cookies.get("user-lang")
  if (language != undefined)
    return language

  if (navigator.languages != undefined)
    language = navigator.languages[0]
  else
    language = navigator.language || window.navigator.userLanguage

  var localLang = language.split("-")
  if (2 == localLang.length)
    language = localLang[0]

  return language
}

function setUserLang(lang) {
  Cookies.set("user-lang", lang)
}

function getMeta(metaName) {
  const metas = document.getElementsByTagName('meta')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('name') === metaName) {
      return metas[i].getAttribute('content')
    }
  }

  return ''
}

function getLink(linkRel) {
  const metas = document.getElementsByTagName('link')

  for (let i = 0; i < metas.length; i++) {
    if (metas[i].getAttribute('rel') === linkRel) {
      return metas[i].getAttribute('href')
    }
  }

  return ''
}

function setI18nStrings(lang, strings) {
  localStorage.setItem('strings-' + lang, JSON.stringify(strings));
}

function getI18nStrings(lang) {
  let strings = localStorage.getItem('strings-' + lang);
  return JSON.parse(strings)
}

function removeI18nStrings(lang) {
  localStorage.removeItem('strings-' + lang);
}

// app.js
// Created at 03/09/2021

let userLang = {
  code: "en",
  name: "English",
}

languages.forEach(element => {
  if (element.code == getUserLang()) {
    userLang = element
    return
  }
});

var defaultLang = {}
const i18n = Vue.observable({
  ...defaultLang,
})

let i18ns = new Map();
i18ns.set(defaultLang.language, defaultLang)

axios.defaults.headers.common['Authorization'] = Cookies.get("access_token")

const plugin = {
  install: function (Vue, options) {
    Vue.prototype.user = {
      lang: userLang,
      logined: logined,
      words: [],
      archives: [],
    }
    Vue.prototype.languages = languages
    Vue.prototype.i18n = i18n
    Vue.prototype.getUserWords = function () {
      axios({
        method: "GET",
        url: queryRestful("/v1/words"),
      }).then(function (resp) {
        Vue.prototype.user.words = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    }
    Vue.prototype.getArchives = function () {
      axios({
        method: "GET",
        url: queryRestful("/v1/archives"),
      }).then(function (resp) {
        Vue.prototype.user.archives = resp.data
      }).catch(function (resp) {
        console.log(resp)
      })
    }

    Vue.prototype.md2html = function (md) {
      var converter = new showdown.Converter({
        'disableForced4SpacesIndentedSublists': 'true',
        'tasklists': 'true',
        'tables': 'true',
        'extensions': ['video', 'audio', 'catalog', 'anchor']
      })
      // KaTeX: math regex: /\$\$([^$]+)\$\$/gm

      return converter.makeHtml(md);
    }
  }
}

const About = { template: fgm_about }
const My = { template: fgm_my }
const Login = { template: fgm_login }

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', component: Index },
    { path: '/p/:id/:code?', name: "article", component: Article, props: true },
    { path: '/drafts', name: 'drafts', component: Drafts },
    { path: '/archive', component: Archive, children: archive_router },
    { path: '/about', name: 'about', component: About },
    { path: '/my', name: 'my', component: My },
    { path: '/settings/personalize', name: 'personalize', component: Personalize },
    { path: '/new/article', name: 'newArticle', component: NewArticle },
    { path: '/d/:id/edit', name: 'editDraft', component: EditDraft, props: true },
    { path: '/d/:id/history', name: 'draftHistories', component: DraftHistories, props: true },
    { path: '/d/:id/history/:hid', name: 'draftHistory', component: DraftHistory, props: true },
    { path: '/d/:id/publish', name: 'publishArticle', component: PublishArticle, props: true },
    { path: '/login', name: 'login', component: Login },
  ]
})

var app = new Vue({
  data: {
    vote: {
      has: false,
      ras: null,
    },
    profilePicture: getLink("icon"),
  },
  router,
  template: logined ? app_home : app_index,

  created() {
    var _this = this
    if (logined) {
      axios({
        method: "GET",
        url: queryRestful("/v1/vote"),
      }).then(function (resp) {
        if (200 == resp.status) {
          _this.vote.ras = resp.data
          _this.vote.has = true
        }
      }).catch(function (resp) {
        console.log(resp)
      })
    }
  },

  methods: {
    onAllow() {
      if (!this.ras) return

      this.__postVote("allowed")
    },
    onRejecte() {
      if (!this.ras) return

      this.__postVote("rejected")
    },
    onGitHubOAuth: function () {
      const github_client_id = getMeta("github_client_id")
      const state = Math.random().toString(36).slice(2)
      const githubOAuthAPI = "https://github.com/login/oauth/authorize?client_id=" + github_client_id + "&state=" + state
      Cookies.set('github_oauth_state', state)
      window.open(githubOAuthAPI, "_self")
    },
    onSignout: function () {
      const githubOAuthAPI = "/signout"
      window.open(githubOAuthAPI, "_self")
    },
    __postVote(status) {
      var _this = this
      axios({
        method: "POST",
        url: queryRestful("/v1/vote"),
        data: {
          id: this.ras.id,
          status: status,
        },
      }).then(function (resp) {
        _this.ras = null
      }).catch(function (resp) {
        console.log(resp)
      })
    },
  }
})

var langCode = getUserLang()
var url = queryStatic("/strings/strings-" + langCode + ".json")

axios({
  method: "GET",
  url: url,
}).then(function (resp) {
  i18ns.set(langCode, resp.data)
  setI18nStrings(langCode, resp.data)
  Object.assign(i18n, resp.data)

  Vue.use(plugin)
  app.$mount('#application--wrap')

  app.getUserWords()
  app.getArchives()
}).catch(function (resp) {
  console.log(resp)
})