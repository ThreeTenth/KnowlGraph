// analytics.js

var city = ""
var timezone = ""
var screenWidth = window.screen.width
var screenHeight = window.screen.height
var parser = new UAParser().setUA(navigator.userAgent)
var device = parser.getResult().device
if (device.model == undefined) {
  device.model = ""
}

axios({ method: "GET", url: "https://ipinfo.io/json", })
  .then((resp) => {
    city = resp.data.city
    timezone = resp.data.timezone
  })

function postAnalyticsPageView(page, title, referrer) {
  axios({
    method: "POST",
    url: queryRestful("/v1/analytics"),
    data: {
      event: "page_view",
      category: page,
      label: title,
      city: city,
      timezone: timezone,
      referrer: referrer,
      url: location.href,
      path: location.path,
      platform: "web",
      platform_version: VersionName,
      device: device.model,
      display: screenWidth + "x" + screenHeight,
      start_time: new Date(),
      version: 1,
    },
  }).then((resp) => {
    console.log(resp.data)
  })
}