// analytics.js

var city = Cookies.get("city")
var timezone = Cookies.get("timezone")
var screenWidth = window.screen.width
var screenHeight = window.screen.height
var parser = new UAParser().setUA(navigator.userAgent)
var device = parser.getResult().device
if (device.model == undefined) {
  device.model = ""
}
var analyticsCode = Cookies.get("analytics_code")
if (analyticsCode === undefined) {
  analyticsCode = getRandomString() + "." + getRandomString()
  Cookies.set("analytics_code", analyticsCode, 30)
}

if (city === undefined || city === null)
  axios({ method: "GET", url: "https://ipinfo.io/json", })
    .then((resp) => {
      city = resp.data.city
      timezone = resp.data.timezone
      Cookies.set("city", city, { expires: 30 })
      Cookies.set("timezone", timezone, { expires: 30 })
    }).catch((err) => {
      Cookies.set("city", "", { expires: 1 })
      Cookies.set("timezone", "", { expires: 1 })
    })

function getAnalyticsUser(form = undefined, success, failure) {
  var copyForm = form ? form : {}
  copyForm.event = "page_view"
  copyForm.group_by = "event"
  copyForm.count_field = "code"
  axios({
    method: "GET",
    url: queryRestful("/v1/analytics", copyForm),
  }).then(success).catch(failure)
}

function getAnalyticsPageView(form = undefined, success, failure) {
  var copyForm = form ? form : {}
  copyForm.event = "page_view"
  axios({
    method: "GET",
    url: queryRestful("/v1/analytics", copyForm),
  }).then(success).catch(failure)
}

function postAnalyticsPageView(page, title, referrer) {
  axios({
    method: "POST",
    url: queryRestful("/v1/analytics"),
    data: {
      code: analyticsCode,
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
    },
  })
}