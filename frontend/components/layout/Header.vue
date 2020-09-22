<template>
  <v-row class="infini-nav" align="center" justify="space-between">
    <v-col sm="6" md="5" lg="3" xl="3">
      <a @click="toggleCollapsed" class="menu-control">
        <a-icon :type="menu ? 'menu-unfold' : 'menu-fold'" />
      </a>
      <span class="logo">infinimesh</span></v-col
    >
    <v-col class="d-none d-lg-block" lg="1" xl="1">
      <v-row>
        <v-col class="nav-button" cols="5" @click="$router.go(-1)">
          <a-icon type="left" />
        </v-col>
        <v-col class="nav-button" cols="5" :offset="14" @click="$router.go(1)">
          <a-icon type="right" />
        </v-col>
      </v-row>
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
  computed: {
    menu: {
      get() {
        return this.$store.getters["window/menu"];
      },
      set(val) {
        this.$store.dispatch("window/toggleMenu", val);
      }
    },
    namespace: {
      get() {
        return this.$store.state.devices.namespace;
      },
      set(val) {
        this.$store.dispatch("devices/setNamespace", val);
      }
    },
    namespaces: {
      deep: true,
      get() {
        return this.$store.state.devices.namespaces;
      }
    }
  },
  mounted() {
    this.namespace = this.$store.state.auth.user.default_namespace.id;
  },
  methods: {
    toggleCollapsed() {
      this.menu = !this.menu;
    }
  }
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
</style>

<style lang="less" scoped>
.nav-button:hover {
  cursor: pointer;
  opacity: 0.5;
  -webkit-filter: grayscale(100%) sepia(100%);
  filter: grayscale(100%) sepia(100%);
}
.infini-nav {
  max-height: 64px;
}
.infini-nav .menu-control {
  font-size: @font-size-xl !important;
}
</style>
