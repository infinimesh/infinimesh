<template>
  <a-menu mode="vertical" v-model="route" id="menu">
    <a-menu-item key="dashboard">
      <a-icon type="cloud" />
      <span>Device Registry</span>
    </a-menu-item>
    <a-sub-menu key="user">
      <span slot="title">
        <a-icon type="user" />
        <span>
          <b>{{ user.name }}</b>
        </span>
      </span>
      <a-menu-item key="logout">
        <a @click="$store.dispatch('logout')">Log Out</a>
      </a-menu-item>
    </a-sub-menu>
  </a-menu>
</template>

<script>
export default {
  computed: {
    user() {
      return this.$store.getters.loggedInUser;
    },
    route: {
      get() {
        return [this.$route.name];
      },
      set(val) {
        this.$router.push({ name: val[0] });
      },
    },
  },
};
</script>

<style scoped>
#menu {
  height: 100%;
}
</style>
