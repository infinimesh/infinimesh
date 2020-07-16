<template>
  <a-layout id="dashboard">
    <a-layout-header class="wide-header">
      <Header v-model="menu_collapsed" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model="menu_collapsed"><Sider /></a-layout-sider>
      <a-layout>
        <a-layout-content>
          <nuxt-child />
        </a-layout-content>
        <a-layout-footer>
          ©2020 — <strong> infinimesh, inc </strong>
          - source code at
          <a href="https://www.github.com/infinimesh/infinimesh" target="_new"
            ><strong style="color: white;">GitHub</strong></a
          >
        </a-layout-footer>
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
  mounted() {
    this.$store.dispatch("devices/getNamespaces");
  }
};
</script>

<style>
.wide-header {
  padding: 0 !important;
}
</style>
