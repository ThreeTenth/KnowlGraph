// vue_time.js

Vue.component('the-time', {
  props: {
    datetime: String,
    full: Boolean,
  },

  computed: {
    reminder: function () {
      let lang = this.i18n.__code
      var options = { weekday: "long", year: "numeric", month: "short", day: "numeric", hour: "numeric", minute: "numeric", hour12: false };
      return new Date(this.datetime).toLocaleString(lang, options)
    },
    fulltime: function () {
      let now = new Date()
      let local = new Date(this.datetime)
      let lang = this.i18n.__code

      var options = { month: "short", day: "numeric", hour: "numeric", minute: "numeric", hour12: false };
      if (now.getFullYear() != local.getFullYear()) {
        options.year = "numeric"
      }
      return new Date(this.datetime).toLocaleString(lang, options)
    },
    timetonow: function () {
      let now = new Date()
      let local = new Date(this.datetime)
      let timer = (now - local) / 1000
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
      } else if (now.getFullYear() == local.getFullYear()) {
        let lang = this.i18n.__code
        var options = { month: "short", day: "numeric" };
        tip = local.toLocaleString(lang, options)
      } else {
        let lang = this.i18n.__code
        var options = { year: "numeric", month: "short", day: "numeric" };
        tip = local.toLocaleString(lang, options)
      }
      return tip
    },
  },

  template: com_the_time,
})