<template>
  <a-menu theme="dark" mode="inline" v-model="route">
    <a-menu-item key="dashboard">
      <a-icon type="cloud" />
      <span>Device Registry</span>
    </a-menu-item>
  </a-menu>
</template>

<script>
export default {
  computed: {
    route: {
      get() {
        return [this.$route.name];
      }
    }
  }
};
</script>
