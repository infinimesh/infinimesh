<template>
  <v-app id="dashboard">
    <v-app-bar app clipped-left>
      <Header />
    </v-app-bar>
    <Sider />
    <v-main>
      <v-container>
        <nuxt-child />
      </v-container>
    </v-main>
    <v-footer app>
      <infinimesh-footer />
    </v-footer>
  </v-app>
</template>

<script>
import Header from "@/components/layout/Header";
import Sider from "@/components/layout/Sider";
import InfinimeshFooter from "@/components/generic/footer.vue";

export default {
  name: "page-dashboard",
  components: {
    Header,
    Sider,
    InfinimeshFooter,
  },
  mounted() {
    this.$store.dispatch("devices/getNamespaces");
  },
};
</script>
