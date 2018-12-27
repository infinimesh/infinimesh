export const APIMixins = {
  methods: {
    getRemoteDevice() {
      this.$http
        .get("http://localhost:8081/devices/" + this.id)
        .then(response => {
          this.$store.dispatch("addDevice", response.body.device);
          this.device = this.$store.getters.getDevice(this.id);
          this.checkbox = this.$store.getters.getDevice(this.id);
        })
        .catch(e => {
          console.log(e);
        });
    }
  }
};
