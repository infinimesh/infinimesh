<template>
  <v-app>
    <Header class="my-1"></Header>
    <v-layout row>
      <v-flex class="mx-4 my-3">
        <router-view></router-view>
      </v-flex>
    </v-layout>
    <Footer></Footer>
  </v-app>
</template>

<script>
import Footer from "./components/Footer.vue";
import Header from "./components/Header.vue";

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
