<template>
  <div id="404"></div>
</template>
<script>
export default {
  mounted() {
    this.$router.push({ name: "dashboard-devices" });
  },
};
</script>
