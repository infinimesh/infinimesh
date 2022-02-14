<template>
  <div id="dashboard_index"></div>
</template>

<script>
export default {
  name: "page-dashboard-index",
  mounted() {
    this.$router.push({ name: "dashboard-devices" });
  },
};
</script>
