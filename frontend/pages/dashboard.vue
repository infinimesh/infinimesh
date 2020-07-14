<template>
  <a-layout>
    <a-layout-header class="wide-header">
      <Header v-model="menu_collapsed" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model="menu_collapsed"><Sider /></a-layout-sider>
      <a-layout>
        <a-layout-content>
          <nuxt-child />
        </a-layout-content>
        <a-layout-footer>Footer</a-layout-footer>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script>
import Header from "@/components/Header";
import Sider from "@/components/Sider";

export default {
  components: {
    Header,
    Sider
  },
  data() {
    return {
      menu_collapsed: false
    };
  },
  fetch({ store, params }) {
    store.dispatch("getNamespaces");
  }
};
</script>

<style>
.wide-header {
  padding: 0 !important;
}
</style>
