<template>
  <v-app id="dashboard">
    <v-app-bar app color="#104e83" clipped-left>
      <Header />
    </v-app-bar>
    <Sider />
    <a-layout class="layout-content">
      <a-layout>
        <a-layout-content
          :style="{
            marginLeft: '80px',
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
