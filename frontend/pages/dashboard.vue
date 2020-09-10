<template>
  <a-layout id="dashboard">
    <a-layout-header class="wide-header">
      <Header />
    </a-layout-header>
    <a-layout class="layout-content">
      <a-layout-sider
        v-model="menu_collapsed"
        :style="{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          zIndex: 99
        }"
      >
        <Sider />
      </a-layout-sider>
      <a-layout>
        <a-layout-content
          :style="{
            marginLeft: menu_collapsed ? '80px' : '200px',
            paddingBottom: '20rem'
          }"
        >
          <nuxt-child />
        </a-layout-content>
        <infinimesh-footer />
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script>
import Header from "@/components/layout/Header";
import Sider from "@/components/layout/Sider";
import InfinimeshFooter from "@/components/generic/footer.vue";

export default {
  components: {
    Header,
    Sider,
    InfinimeshFooter,
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
      },
    },
  },
};
</script>

<style>
.wide-header {
  padding: 0 !important;
  position: fixed;
  z-index: 1;
  width: 100%;
}
</style>
<style lang="less" scoped>
.layout-content {
  margin-top: @layout-header-height;
}
</style>
