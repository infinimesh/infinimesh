<template>
  <v-navigation-drawer
      stateless
      value="true"
    >
      <v-list>
        <div
          v-for="(page, i) of pages"
          :key="i"
        >
          <v-list-tile
            :to="{ name: page[1] }"
            exact
          >
            <v-list-tile-action />
            <v-list-tile-title
              class="body-2"
            >
              {{ page[0] }}
            </v-list-tile-title>
          </v-list-tile>
          <v-divider></v-divider>
        </div>
      </v-list>
    </v-navigation-drawer>
</template>

<script>
export default {
  data: () => ({
    pages: [
      ["Devices", "View devices"],
      ["Assets", "Manage devices"],
      ["Monitoring", "View devices"],
      ["Integrations", "View devices"],
      ["Accounts", "Account Management"]
    ]
  })
};
</script>

<style lang="css" scoped>
</style>
