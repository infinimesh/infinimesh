<template>
  <a-menu mode="inline" v-model="route" id="menu">
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
      },
      set(val) {
        this.$router.push({ name: val[0] });
      }
    }
  }
};
</script>

<style scoped>
#menu {
  height: 100%;
}
</style>
