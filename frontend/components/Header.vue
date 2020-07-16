<template>
  <a-row class="gay-theme-nav">
    <a-col :xs="{ span: 1, offset: 1 }" :md="{ span: 1, offset: 0 }">
      <a @click="toggleCollapsed" class="menu-control">
        <a-icon :type="value ? 'menu-unfold' : 'menu-fold'" />
      </a>
    </a-col>

    <a-col :xs="{ span: 12, offset: 2 }" :sm="{ offset: 1 }" :md="{ span: 8 }">
      <a-row type="flex" justify="start">
        <a-col>
          <img
            src="@/assets/infinimesh_logo.png"
            alt="infinimesh Logo"
            class="logo"
          />
        </a-col>
        <a-col class="logo">infinimesh<span>.io</span></a-col>
      </a-row>
    </a-col>
    <a-col
      :xs="{ span: 6, offset: 1 }"
      :sm="{ span: 4, offset: 4 }"
      :md="{ span: 4, offset: 4 }"
      :lg="{ span: 4, offset: 4 }"
      :xl="{ span: 6, offset: 4 }"
      :xxl="{ span: 3, offset: 7 }"
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
    <a-col
      :xs="{ span: 0 }"
      :md="{ span: 5, offset: 1 }"
      :xl="{ span: 3, offset: 1 }"
    >
      <a-row type="flex" justify="end">
        <a-col :span="8">
          <a-avatar>R</a-avatar>
        </a-col>
        <a-col
          class="user"
          :xs="{ span: 0 }"
          :sm="{ span: 0 }"
          :md="{ span: 16 }"
        >
          root
        </a-col>
      </a-row>
    </a-col>
  </a-row>
</template>

<script>
export default {
  props: ["value"],
  computed: {
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
      this.value = !this.value;
      this.$emit("input", this.value);
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
