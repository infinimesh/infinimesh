<template>
  <v-card
  color="grey lighten-4"
  flat
  >
  <v-toolbar color="primary lighten-1" dark>
    <v-toolbar-title
      class="headline text-uppercase white--text"
    >
      <span
        @click="$router.push('/devices')"
        style="cursor: pointer"
      >
        Infinimesh
      </span>
    </v-toolbar-title>
    <v-spacer></v-spacer>
    <v-menu :nudge-width="100" offset-y>
      <v-toolbar-title slot="activator">
        <span
        color="white"
        >Devices</span>
        <v-icon
        dark
        class="mr-5"
        >arrow_drop_down</v-icon>
      </v-toolbar-title>
      <v-list>
        <v-list-tile
          v-for="page in deviceLinks"
          :key="page"
          :to="{ name: page }"
          exact
          >
          <v-list-tile-title
            v-text="page"
          ></v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
    <v-menu
      :nudge-width="100"
      offset-y
    >
      <v-toolbar-title
      slot="activator"
      >
        <v-icon
        large
        color="white"
        >
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
          ></v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
    </v-btn>
  </v-toolbar>
</v-card>
</template>

<script>
export default {
  data() {
    return {
      deviceLinks: ["View devices", "Manage devices", "Register device"],
      accountLinks: ["Logout"]
    };
  }
};
</script>

<style>
</style>
