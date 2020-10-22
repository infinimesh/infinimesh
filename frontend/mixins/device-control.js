export default {
  methods: {
    refresh() {
      this.$axios
        .get(`/api/devices/${this.device.id}`)
        .then(res => {
          this.device = res.data.device;
          this.socket = new WebSocket(
            `wss://${this.$config.baseURL.replace("https://", "")}/devices/${
              this.device.id
            }/state/stream`,
            ["Bearer", this.$auth.getToken("local").split(" ")[1]]
          );
          this.socket.onmessage = msg => {
            let response = JSON.parse(msg.data).result;
            if (response)
              this.device.state.shadow.reported = response.reportedState;
          };
          window.addEventListener("beforeunload", function(event) {
            this.socket.close();
          });
        })
        .catch(res => {
          if (res.response.status == 404) {
            this.$notification.error({
              message: "Device wasn't found",
              description: "Redirecting..."
            });
          } else if (res.response.status == 403) {
            this.$notification.error({
              message: "You have no access to this device",
              description: "Redirecting..."
            });
          }
          this.$router.push({ name: "dashboard-devices" });
        });
      this.deviceStateGet();
    },
    /**
     * Obtains device state(s) (desired and reported) and merges them into deviceObject
     */
    async deviceStateGet() {
      await this.$axios
        .get(`/api/devices/${this.device.id}/state`)
        .then(res => {
          this.device = {
            ...this.device,
            state: res.data
          };
        });
    },
    /**
     * Performs PATCH /device/id and changes desired state to given.
     * Used by device-state component so it invokes callback after patch response handling is done
     * @param {String} state, the desired state as in JSON String
     */
    handleStateUpdate(state, callback) {
      this.$axios({
        url: `/api/devices/${this.device.id}/state`,
        method: "patch",
        data: state
      })
        .then(res => {
          this.deviceStateGet();
        })
        .catch(res => {
          console.error(res);
        })
        .then(() => {
          callback();
        });
    },
    handleDeviceDelete() {
      this.$axios({
        url: `/api/devices/${this.device.id}`,
        method: "delete"
      })
        .then(() => {
          this.$message.success("Device successfuly deleted!");
          this.$store.dispatch("devices/get");
          this.$router.push({ name: "dashboard-devices" });
        })
        .catch(e => {
          this.$notification.error({
            message: "Error deleting device",
            description: e.response.data.message
          });
        });
    },
    handleToogleDevice(refresh = true) {
      this.handleDeviceUpdate(
        {
          enabled: !this.device.enabled
        },
        {
          refresh: refresh,
          success: () => {
            this.$message.success(
              `Device successfuly ${
                this.device.enabled ? "disabled" : "enabled"
              }!`
            );
          },
          error: () => {
            this.$notification.error({
              message: `Error ${
                device.enabled ? "disabling" : "enabling"
              } device`,
              description: e.response.data.message
            });
          }
        }
      );
    },
    handleDeviceUpdate(data, { refresh, success, error }) {
      this.$axios({
        url: `/api/devices/${this.device.id}`,
        method: "patch",
        data: data
      })
        .then(() => {
          if (success) success();
        })
        .catch(e => {
          if (error) error(e);
        })
        .then(() => {
          if (refresh) this.refresh();
          this.$store.dispatch("devices/get");
        });
    }
  }
};
