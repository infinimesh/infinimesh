export const APIMixins = {
  methods: {
    getRemoteDevice() {
      this.$http
        .get("devices/" + this.id)
        .then(response => {
          this.$store.dispatch("addDevice", response.body.device);
          this.device = this.$store.getters.getDevice(this.id);
          this.checkbox = this.$store.getters.getDevice(this.id).enabled;
        })
        .catch(e => {
          console.log(e);
        });
    }
  }
};

export const storeRemoteDevices = {
  methods: {
    storeRemoteDevices() {
      this.$http
        .get("devices")
        .then(response => {
          let device = {};
          for (device of response.body.devices) {
            this.$store.dispatch("addDevice", device);
          }
        })
        .catch(e => {
          console.log(e);
        });
    }
  }
};
