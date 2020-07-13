<template>
  <a-row class="gay-theme-nav">
    <a-col :span="1">
      <a @click="toggleCollapsed" class="menu-control">
        <a-icon :type="value ? 'menu-unfold' : 'menu-fold'" />
      </a>
    </a-col>

    <a-col :span="6">
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
    <a-col :span="3" :offset="10">
      <a-select
        style="width: 100%"
        :default-value="$store.state.devices.namespace"
      >
        <a-select-option :key="ns.id" :value="ns.id" v-for="ns in namespaces"
          >NS: {{ ns.name }}</a-select-option
        >
      </a-select>
    </a-col>
  </a-row>
</template>

<script>
export default {
  props: ["value"],
  computed: {
    namespaces() {
      return this.$store.state.devices.namespaces;
    }
  },
  fetch({ store }) {
    store.commit(
      "devices/namespace",
      $store.state.auth.user.default_namespace.id
    );
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
  font-size: 1.75rem;
  font-weight: 500;
  white-space: nowrap;
  font-family: Exo;
  color: white;
}
.root {
  box-shadow: 0 8px 20px 0 rgba(40, 37, 89, 0.6);
}
</style>
