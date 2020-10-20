<template>
  <div>
    <nuxt />
  </div>
</template>

<script>
export default {
  mounted() {
    console.log(
      "%c infinimesh.io %c v0.1.5 %c",
      "background:#104e83; padding: 1px; border-radius: 3px 0 0 3px; color: #fff; font-size: 18px; font-weight: 500;",
      "background:#35495e; padding: 1px; border-radius: 0 3px 3px 0; color: #fff; font-size: 18px;",
      "background:transparent"
    );
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
