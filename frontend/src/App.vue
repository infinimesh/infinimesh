<template>
  <v-app id="inspire">
    <Header class="my-1"></Header>
    <v-layout row>
      <v-flex>
        <Sidebar></Sidebar>
      </v-flex>
      <v-flex xs9 class="mx-4 my-3">
        <router-view></router-view>
      </v-flex>
    </v-layout>
    <Footer></Footer>
  </v-app>
</template>

<script>
import Footer from "./components/Footer.vue";
import Header from "./components/Header.vue";
import Sidebar from "./components/Sidebar.vue";
import RegisterDevice from "./views/RegisterDevice.vue";

export default {
  name: "App",
  components: {
    Footer,
    Header,
    Sidebar,
    RegisterDevice
  },
  data() {
    return {
      //
    };
  }
};
</script>
