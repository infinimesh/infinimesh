<template>
  <v-app>
    <Header class="my-1"></Header>
        <router-view></router-view>
    <Footer></Footer>
  </v-app>
</template>

<script>
import Footer from "./components/TheFooter.vue";
import Header from "./components/TheHeader.vue";

export default {
  name: "App",
  components: {
    Footer,
    Header
  },
  data() {
    return {
      //
    };
  }
};
</script>
