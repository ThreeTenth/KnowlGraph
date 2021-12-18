// vue_analytics-chart.js

const color = ["#7FDBFF", "#F012BE", "#39CCCC", "#B10DC9", "#01FF70", "#85144B", "#3D9970", "#1a5cff"]

Vue.component('analytics-chart', {
  props: {
    groupBy: {
      type: String,
    },
    countField: String,
  },
  data: function () {
    return {
      pageStatus: 0,
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
            let groupByName = this.groupBy ?
              (
                v[this.groupBy] ?
                  v[this.groupBy] :
                  "Unknown"
              ) :
              "浏览量"

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
      new Chart(ctx, {
        type: 'bar',
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
              display: true,
              text: this.countField === "code" ? "用户数" : '浏览量'
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
    var now = new Date()
    var _1day = new Date(now.getTime() - 24 * 60 * 60 * 1000)
    var _2day = new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000)
    var _7day = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)

    var form = {
      time_start: [
        _7day.toJSON(),
        _2day.toJSON(),
        _1day.toJSON()
      ],
      time_unit: 1440,
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