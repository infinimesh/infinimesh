<template>
  <a-row class="gay-theme-nav" type="flex" align="middle">
    <a-col :xs="{ span: 1, offset: 1 }" :md="{ span: 1, offset: 0 }">
      <a-row type="flex" justify="center">
        <a @click="toggleCollapsed" class="menu-control">
          <a-icon :type="menu ? 'menu-unfold' : 'menu-fold'" />
        </a>
      </a-row>
    </a-col>

    <a-col
      :xs="{ span: 12, offset: 2 }"
      :sm="{ span: 12, offset: 1 }"
      :md="{ span: 8 }"
      :xxl="{ span: 7 }"
      class="logo"
      >infinimesh.io</a-col
    >
    <a-col
      :xs="{ span: 0 }"
      :sm="{ span: 0 }"
      :md="{ span: 1, offset: 2 }"
      :lg="{ span: 1, offset: 2 }"
      :xl="{ span: 1, offset: 2 }"
      :xxl="{ span: 1, offset: 4 }"
    >
      <a-row>
        <a-col class="nav-button" :span="5" @click="$router.go(-1)">
          <a-icon type="left" />
        </a-col>
        <a-col class="nav-button" :span="5" :offset="14" @click="$router.go(1)">
          <a-icon type="right" />
        </a-col>
      </a-row>
    </a-col>
    <a-col
      :xs="{ span: 6, offset: 1 }"
      :sm="{ span: 4, offset: 4 }"
      :md="{ span: 6, offset: 4 }"
      :lg="{ span: 6, offset: 3 }"
      :xl="{ span: 6, offset: 4 }"
      :xxl="{ span: 4, offset: 5 }"
    >
      <a-select style="width: 100%" v-model="namespace">
        <a-select-option
          :key="ns.id"
          :value="ns.id"
          :label="ns.name"
          v-for="ns in namespaces"
          >NS: {{ ns.name }}</a-select-option
        >
      </a-select>
    </a-col>
  </a-row>
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
  background-color: @primary-color-dark;
}
.gay-theme-nav {
  max-height: 64px;
}
.gay-theme-nav .menu-control {
  font-size: @font-size-xl !important;
}
</style>
