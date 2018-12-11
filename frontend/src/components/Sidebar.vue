<template>
    <v-navigation-drawer
      stateless
      value="true"
      class="grey lighten-4"
    >
      <v-list>
        <v-list-group
            no-action
            sub-group
            value="true"
            v-for="(category, i) in Object.keys(pages)"
            :key="i"
          >
            <v-list-tile slot="activator">
              <v-list-tile-title class="title" v-text="category"></v-list-tile-title>
            </v-list-tile>

            <v-list-tile
              v-for="(page, i) in pages[category]"
              :key="i"
              @click="navigateTo(page)"
            >
              <v-list-tile-title v-text="page"></v-list-tile-title>
              <v-list-tile-action></v-list-tile-action>
            </v-list-tile>
          </v-list-group>
      </v-list>
    </v-navigation-drawer>
</template>

<script>
import router from "../router";

export default {
  data: () => ({
    pages: {
      Devices: ["View your devices", "Register device", "Delete device"],
      Documentation: ["API", "Examples", "Tutorials"],
      Settings: ["Account", "Profile"]
    }
  }),
  methods: {
    navigateTo(page) {
      this.$router.push({ name: page });
    }
  }
};
</script>

<style>
</style>
