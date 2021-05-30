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