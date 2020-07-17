<template>
  <a-layout id="dashboard">
    <a-layout-header class="wide-header">
      <Header />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model="menu_collapsed"><Sider /></a-layout-sider>
      <a-layout>
        <a-layout-content>
          <nuxt-child />
        </a-layout-content>
        <a-layout-footer>
          <a-row type="flex" justify="center">
            <a-col :xs="24" :sm="18" :md="12" :lg="10" :xl="8">
              ©2020 — <strong> infinimesh, inc </strong>
              - source code at
              <a
                href="https://www.github.com/infinimesh/infinimesh"
                target="_new"
                ><strong style="color: white;">GitHub</strong></a
              >
            </a-col>
          </a-row>
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
  mounted() {
    this.$store.dispatch("devices/getNamespaces");
  },
  computed: {
    menu_collapsed: {
      get() {
        return this.$store.getters["window/menu"];
      },
      set(val) {
        this.$store.dispatch("window/toggleMenu", val);
      }
    }
  }
};
</script>

<style>
.wide-header {
  padding: 0 !important;
}
</style>
