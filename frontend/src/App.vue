<template>
  <v-app>
    <Header></Header>
    <v-layout
      row
      wrap
    >
        <Sidebar />
      <v-flex>
        <router-view></router-view>
      </v-flex>
    </v-layout>
    <Footer></Footer>
  </v-app>
</template>

<script>
import Footer from "./components/TheFooter.vue";
import Header from "./components/TheHeader.vue";
import Sidebar from "./components/TheSidebar.vue";

export default {
  name: "App",
  components: {
    Footer,
    Header,
    Sidebar
  },
  data() {
    return {
      //
    };
  }
};
</script>
