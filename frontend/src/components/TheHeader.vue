<template>
  <v-toolbar  dark color="primary" height="60px">
    <img src="../../images/infinimesh_small_bw.png" alt="Infinimesh" @click="$router.push('/devices')" style="cursor: pointer" /> &nbsp; &nbsp;
    <v-layout>
      <v-toolbar-items>
        <v-flex xs6>
          <namespace-menu light></namespace-menu>
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
