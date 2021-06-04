// personalize.js

const Personalize = {
  computed: {
    currentLang: function() {
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
        Object.assign(i18n, defaultLang)
        return
      } else if (i18ns.has(langCode)) {
        Object.assign(i18n, i18ns.get(langCode))
        return
      } else {
        let strings = getI18nStrings(langCode)
        if (null != strings && strings.__version == defaultLang.__version) {
          i18ns.set(langCode, strings)
          Object.assign(i18n, strings)
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
        Object.assign(i18n, resp.data)
      }).catch(function (resp) {
        _this.user.lang = old
        setUserLang(old.code)
      })
    },
  },

  template: fgm_personalize
}