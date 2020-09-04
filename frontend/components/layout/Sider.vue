<template>
  <a-menu mode="vertical" v-model="route" id="menu">
    <a-menu-item :key="page.link" v-for="page in pagesFiltered">
      <a-icon :type="page.icon" />
      <span>{{page.title}}</span>
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
  data() {
    return {
      pages: [
        { title: "Device Registry", icon: "cloud", link: "dashboard-devices" },
        { title: "Accounts", icon: "idcard", link: "dashboard-accounts" },
        {
          title: "Namespaces",
          icon: "folder-open",
          link: "dashboard-namespaces",
        },
      ],
    };
  },
  computed: {
    pagesFiltered() {
      return this.pages.filter((page) => this.allowedScope(page.link));
    },
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
  methods: {
    allowedScope(scope) {
      return this.$store.getters["window/hasAccess"](scope);
    },
  },
};
</script>

<style scoped>
#menu {
  height: 100%;
}
</style>
