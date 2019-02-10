<template>
  <v-navigation-drawer
      stateless
      value="true"
    >
      <v-list>
        <div
          v-for="page in pages"
          :key="page"
        >
          <v-list-tile>
            <v-list-tile-action />
            <v-list-tile-title>{{ page }}</v-list-tile-title>
          </v-list-tile>
          <v-divider></v-divider>
        </div>
      </v-list>
    </v-navigation-drawer>
</template>

<script>
export default {
  data: () => ({
    pages: ["Devices", "Assets", "Accounts", "Monitoring", "Integrations"]
  })
};
</script>

<style lang="css" scoped>
</style>
