<template>
  <div>
    <h1>Shadow: {{ id }}</h1>
  </div>
</template>

<script>
export default {
  data() {
    return {
      id: this.$route.params.id,
      device: {
        value: false,
        status: "active",
        name: "Device 1",
        location: "Düsseldorf",
        tags: "test"
      }
    };
  }
};
</script>

<style lang="css">
</style>
