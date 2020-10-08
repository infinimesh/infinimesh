<template>
  <div>
    <nuxt />
  </div>
</template>

<script>
export default {
  mounted() {
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
