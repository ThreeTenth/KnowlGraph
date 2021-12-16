// dashboard.js

const Dashboard = {
  data: () => {
    return {
    }
  },
  // computed: {},
  methods: {
    getAnalyticsPageViewSuccess() { },
    getAnalyticsPageViewFailure() { },
  },
  created() {
    getAnalyticsPageView((resp) => { }, (err) => { })
  },
  template: fgm_dashboard,
}