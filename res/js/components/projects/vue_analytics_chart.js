// vue_analytics-chart.js

const color = ["#7FDBFF", "#F012BE", "#39CCCC", "#B10DC9", "#01FF70", "#85144B", "#3D9970", "#1a5cff"]

Vue.component('analytics-chart', {
  props: {
    theme: String,
    groupBy: {
      type: String,
    },
    countField: String,
    title: String,
    timeUnit: Number,
    days: Number,
  },
  data: function () {
    return {
      pageStatus: 0,
      chartType: "bar",
    }
  },
  methods: {
    analytics(arr) {
      var i = 0
      var labels = []
      var datasets = []
      for (const key in arr) {
        if (Object.hasOwnProperty.call(arr, key)) {
          const value = arr[key];
          var data = []
          value.forEach(v => {
            let countFileName = this.countField || "浏览量"
            countFileName = countFileName == "code" ? "用户数" : countFileName

            let groupByName = this.groupBy ?
              (
                v[this.groupBy] ?
                  v[this.groupBy] :
                  countFileName
              ) :
              countFileName

            let idx = labels.indexOf(groupByName)
            if (-1 == idx) {
              labels.push(groupByName)
              idx = labels.length - 1
            }
            data[idx] = v.count
          });
          for (let index = 0; index < labels.length; index++) {
            if (data[index] == undefined) {
              data[index] = 0
            }
          }
          var label = new Date(key).toLocaleDateString()
          datasets.push({
            label: label,
            data: data,
            backgroundColor: color[i++],
          })
        }
      }
      // console.log(labels);
      // console.log(datasets);
      const ctx = this.$refs.chart.getContext('2d');
      const title = this.title
      new Chart(ctx, {
        type: this.chartType,
        data: {
          labels: labels,
          datasets: datasets,
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: 'top',
            },
            title: {
              display: !isEmpty(title),
              text: title
            }
          }
        },
      })
    },
    analyticsPageViewSuccess(resp) {
      this.pageStatus = resp.status
      setTimeout(() => {
        this.analytics(resp.data)
      }, 0);
    },
    analyticsPageViewFailure(err) { },
  },
  created() {
    this.chartType = this.theme || this.chartType
    var _day = 24 * 60 * 60 * 1000
    var today = new Date()
    today.setHours(0)
    today.setMinutes(0)
    today.setSeconds(1)
    var timeStart = [today.toJSON()]
    if (this.days) {
      for (let index = 1; index < this.days; index++) {
        timeStart.push(new Date(today.getTime() - index * _day).toJSON())
      }
    }

    var form = {
      time_start: timeStart,
    }
    if (this.timeUnit) {
      form.time_unit = this.timeUnit
    }
    if (this.groupBy) {
      form.group_by = this.groupBy
    } else {
      form.group_by = "event"
    }
    if (this.countField) {
      form.count_field = this.countField
    }
    // 统计最近两天的页面访问量
    getAnalyticsPageView(
      form,
      this.analyticsPageViewSuccess,
      this.analyticsPageViewFailure)
  },
  template: com_analytics_chart,
})