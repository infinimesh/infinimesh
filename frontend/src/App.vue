<template>
  <v-app>
    <Header />
    <v-layout
      row
      wrap
    >
      <Sidebar
        v-if="showSidebar"
      ></Sidebar>
      <v-flex>
        <v-container
          fluid
          class="pa-0 ma-0"
        >
          <v-card
            flat
            :min-height="windowHeight"
          >
            <router-view />
          </v-card>
        </v-container>
      </v-flex>
    </v-layout>
    <Footer />
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
      showSidebar: false
    };
  },
  computed: {
    // bug: this is not reactive to window.innerHeight
    windowHeight() {
      return window.innerHeight - 128;
    }
  },
  updated() {
    this.showSidebar = localStorage.token;
  },
  created() {
    this.showSidebar = localStorage.token;
  }
};
</script>
