// vue_auth_camera.js

Vue.component('auth-camera', {

  data: function () {
    return {
      codeReader: undefined,
      videoInputDevices: [],
      inputDeviceIndex: 0,

      activateSuccessed: false,
      activateTerminalName: "",
    }
  },

  computed: {
  },

  methods: {

    switchDeviceId() {
      var index = this.inputDeviceIndex
      index = this.videoInputDevices.length <= (index + 1) ? 0 : (index + 1)
      this.codeReader.reset()
      this.startInputDevice(this.videoInputDevices, index)
      this.inputDeviceIndex = index
    },

    startInputDevice(videoInputDevices, index) {
      if (0 == videoInputDevices.length) {
        return
      }

      let selectedDeviceId = videoInputDevices[index].deviceId

      this.codeReader.decodeFromVideoDevice(selectedDeviceId, 'video', (result, err) => {
        if (result) {
          var arr = result.text.split("/")
          var challenge = arr[arr.length - 1]
          postActivateTerminal(challenge, (resp) => {
            if (logined) {
              this.authorizedHandler(resp, challenge)
            } else {
              var state = arr[arr.length - 2]
              this.unAuthorizedHandler(resp, challenge, state)
            }
          }, this.activateTerminalFailure)
        }
        if (err && !(err instanceof ZXing.NotFoundException)) {
          console.error(err)
          this.toast("扫描失败：" + err, "error")
        }
      })
    },

    authorizedHandler(resp, challenge) {
      this.codeReader.reset()
      this.$emit('authorize', resp, challenge)
    },
    unAuthorizedHandler(resp, challenge, state) {
      this.codeReader.reset()
      this.activateTerminalName = resp.data.name
      this.activateSuccessed = true
      this.checkChallengeState(challenge, state)
    },
    activateTerminalFailure(err) {
      this.toast("激活失败：" + err, "error")
    },

    checkChallengeState(challenge, state) {
      this.$emit('challenge', challenge, state)
      var interval = window.setInterval(() => {
        getScanChallenge(challenge, (resp) => {
          if (4 == resp.data.state) {
            window.clearInterval(interval)
          }
          this.$emit('result', resp.data)
        }, function (err) {
          window.clearInterval(interval)
        })
      }, 1000);
    },
  },

  created() {
    const codeReader = new ZXing.BrowserMultiFormatReader()
    codeReader.listVideoInputDevices()
      .then((videoInputDevices) => {
        this.startInputDevice(videoInputDevices, 0)
        this.videoInputDevices = videoInputDevices
        this.inputDeviceIndex = 0
      })
      .catch((err) => {
        console.error(err)
      })

    this.codeReader = codeReader
  },

  beforeDestroy() {
    this.codeReader.reset()
    this.codeReader = null
  },


  template: com_auth_camera,
})