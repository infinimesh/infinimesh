<template>
  <v-card color="grey lighten-4" flat>
    <v-toolbar color="primary lighten-1" dark height="80px">
      <v-toolbar-title class="headline text-uppercase white--text">
        <img
          src="../../images/infinimesh_small_bw.png"
          alt="Infinimesh"
          @click="$router.push('/devices')"
          style="cursor: pointer"
        />
      </v-toolbar-title>
      <namespace-menu />
      <v-spacer></v-spacer>
      <v-menu :nudge-width="100" offset-y left>
        <v-toolbar-title slot="activator">
          <v-icon large color="white">
            account_circle
          </v-icon>
        </v-toolbar-title>
        <v-list>
          <v-list-tile
            v-for="page in accountLinks"
            :key="page"
            :to="{ name: page }"
            exact
          >
            <v-list-tile-title
              v-text="page"
              @click="logout"
            ></v-list-tile-title>
          </v-list-tile>
        </v-list>
      </v-menu>
    </v-toolbar>
  </v-card>
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
