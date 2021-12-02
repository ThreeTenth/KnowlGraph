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
        url: queryStatic("/strings/strings-" + langCode + ".json"),
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

  created() {
    document.title = "人性化 -- KnowlGraph"
  },

  template: fgm_personalize
}