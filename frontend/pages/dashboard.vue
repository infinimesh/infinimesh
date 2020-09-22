<template>
  <v-app id="dashboard">
    <v-app-bar app class="wide-header" color="#104e83">
      <Header />
    </v-app-bar>
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
  </v-app>
</template>

<script>
import Header from "@/components/layout/Header";
import Sider from "@/components/layout/Sider";
import InfinimeshFooter from "@/components/generic/footer.vue";

export default {
  components: {
    Header,
    Sider,
    InfinimeshFooter
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

<style lang="less" scoped>
.layout-content {
  margin-top: @layout-header-height;
}
</style>
