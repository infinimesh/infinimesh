<template>
  <v-row class="infini-nav" align="center">
    <v-col cols="1">
      <v-row justify="center">
        <a @click="toggleCollapsed" class="menu-control">
          <a-icon :type="menu ? 'menu-unfold' : 'menu-fold'" />
        </a>
      </v-row>
    </v-col>

    <v-col sm="6" md="4" xl="3" class="logo">infinimesh.io</v-col>
    <v-col sm="0" md="1" lg="1" xl="1" offset-sm="1" offset-lg="2">
      <v-row>
        <v-col class="nav-button" cols="5" @click="$router.go(-1)">
          <a-icon type="left" />
        </v-col>
        <v-col class="nav-button" cols="5" :offset="14" @click="$router.go(1)">
          <a-icon type="right" />
        </v-col>
      </v-row>
    </v-col>
    <v-col sm="2" md="3" lg="3" xl="3" offset-sm="2" offset="2">
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
img.logo {
  height: 35px;
}
div.logo {
  padding: 0 1.25rem;
  font-weight: 500;
  white-space: nowrap;
  font-family: Exo;
  color: white;
}
@media screen and (max-width: 448px) {
  div.logo {
    font-size: 1.2rem;
    padding: 0;
    padding-left: 0.3rem;
  }
  img.logo {
    height: 30px;
  }
}
@media screen and (min-width: 448px) {
  div.logo {
    font-size: 1.5rem;
  }
}
@media screen and (min-width: 576px) {
  div.logo {
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
