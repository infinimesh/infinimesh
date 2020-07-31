<template>
  <div id="dashboard_index"></div>
</template>

<script>
export default {
  mounted() {
    this.$router.push({ name: "dashboard-devices" });
  }
};
</script>
