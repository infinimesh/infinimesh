<template>
  <v-toolbar  dark color="primary" height="70px">
    <img src="../../public/images/infinimesh_small_bw.png" alt="infinimesh" @click="$router.push('/devices')" style="cursor: pointer" align="top"/> &nbsp; &nbsp; &nbsp; &nbsp;
    <v-layout>
      <v-toolbar-items>
        <v-flex xs6>
          <namespace-menu light align="top"></namespace-menu>
        </v-flex>
        <!-- <div>dd</div> -->
        <!-- <v-select></v-select> -->
      </v-toolbar-items>
  
    </v-layout>
  </v-toolbar>
</template>

<script>
import NamespaceMenu from "./NamespaceMenu.vue";

export default {
  data() {
    return {
      accountLinks: ["Logout"]
    };
  },
  methods: {
    logout() {
      localStorage.token = "";
    }
  },
  components: {
    NamespaceMenu
  }
};
</script>

<style>
</style>
