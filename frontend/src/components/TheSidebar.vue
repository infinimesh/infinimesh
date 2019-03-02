<template>
  <v-navigation-drawer stateless value="true">
    <v-list>
      <div v-for="(page, i) of pages" :key="i">
        <v-list-tile
          :to="{
            name: page[1],
            params: { namespace }
           }"
           exact
          >
          <v-list-tile-action />
          <v-list-tile-title class="body-2">
            {{ page[0] }}
          </v-list-tile-title>
        </v-list-tile>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  data: () => ({
    pages: [
      ["Device Registry", "Devices"],
      ["Nodes", "Device Management"],
      ["Account Management", "Account Management"]
    ]
  }),
  computed: {
    ...mapGetters({ namespace: "getNamespace" })
  }
};
</script>

<style lang="css" scoped></style>
