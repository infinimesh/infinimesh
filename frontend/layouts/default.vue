<template>
  <div>
    <nuxt />
  </div>
</template>

<script>
import axios from "axios";

export default {
  mounted() {
    axios
      .get("https://api.github.com/repos/InfiniteDevices/infinimesh/releases")
      .then((res) => {
        this.$store.dispatch("window/setVersion", res.data[0]);
      })
      .catch(() => {
        console.error("Can't get release tag from GitHub");
      });

    this.$nextTick(() => {
      window.addEventListener("resize", this.onWindowResize);
      this.onWindowResize();
    });
    this.$notification.config({
      placement: "bottomRight",
    });
  },
  beforeDestroy() {
    window.removeEventListener("resize", this.onWindowResize);
  },
  methods: {
    onWindowResize() {
      this.$store.dispatch("window/set", {
        height: window.innerHeight,
        width: window.innerWidth,
      });
    },
  },
};
</script>

<style>
html {
  font-family: Exo, "Source Sans Pro", -apple-system, BlinkMacSystemFont,
    "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}
</style>
