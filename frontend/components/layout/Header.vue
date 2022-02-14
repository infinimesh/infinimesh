<template>
  <v-row class="infini-nav" align="center" justify="space-between">
    <v-col sm="6" md="5" lg="3" xl="3">
      <span class="logo"
        >infinimesh{{
          $colorMode.preference.includes("hallo") ? " 🎃" : ""
        }}</span
      >
    </v-col>
    <v-col cols="1">
      <a-row type="flex" justify="space-between" align="middle" :gutter="5">
        <a-col class="nav-button" @click="$router.go(-1)" v-if="show_nav_btns">
          <a-icon type="left" />
        </a-col>
        <a-col
          v-if="top_action"
          class="nav-button"
          @click="top_action.callback"
        >
          <a-icon :type="top_action.icon" style="font-size: 1.3rem" />
        </a-col>
        <a-col class="nav-button" @click="$router.go(1)" v-if="show_nav_btns">
          <a-icon type="right" />
        </a-col>
      </a-row>
    </v-col>
    <v-col sm="5" md="5" lg="3" xl="3">
      <a-select style="width: 100%" v-model="namespace">
        <a-select-option
          :key="ns.id"
          :value="ns.id"
          :label="ns.name"
          v-for="ns in namespaces"
          >NS: {{ ns.name }}</a-select-option
        >
      </a-select>
    </v-col>
  </v-row>
</template>

<script>
export default {
  name: "layout-header",
  computed: {
    menu: {
      get() {
        return this.$store.getters["window/menu"];
      },
      set(val) {
        this.$store.dispatch("window/toggleMenu", val);
      },
    },
    show_nav_btns() {
      return ["lg", "xl", "xxl"].includes(this.$store.state.window.gridSize);
    },
    top_action() {
      return this.$store.getters["window/topAction"];
    },
    namespace: {
      get() {
        return this.$store.state.devices.namespace;
      },
      set(val) {
        let old = this.namespace;
        this.$store.dispatch("devices/setNamespace", val);
        if (this.$route.name != "dashboard-devices" && old != "") {
          this.$router.push({
            name: "dashboard-devices",
          });
        }
      },
    },
    namespaces: {
      deep: true,
      get() {
        return this.$store.state.devices.namespaces;
      },
    },
  },
  mounted() {
    this.namespace = this.$store.state.auth.user.default_namespace.id;
  },
  methods: {
    toggleCollapsed() {
      this.menu = !this.menu;
    },
  },
};
</script>

<style scoped>
span.logo {
  padding: 0 1.25rem;
  font-weight: 500;
  white-space: nowrap;
  font-family: Exo;
  color: white;
}
@media screen and (max-width: 448px) {
  span.logo {
    font-size: 1.2rem;
    padding: 0;
    padding-left: 0.3rem;
  }
}
@media screen and (min-width: 448px) {
  span.logo {
    font-size: 1.5rem;
  }
}
@media screen and (min-width: 576px) {
  span.logo {
    font-size: 1.75rem;
  }
}
div.user {
  color: white;
  text-align: left;
}

.nav-button {
  color: var(--secondary-color);
}
.nav-button:hover {
  cursor: pointer;
  opacity: 0.5;
  -webkit-filter: grayscale(100%) sepia(100%);
  filter: grayscale(100%) sepia(100%);
}
.infini-nav {
  max-height: 64px;
  background: var(--primary-color);
  padding-right: 12px;
}
.infini-nav .menu-control {
  font-size: var(--font-size-xl) !important;
}
</style>
