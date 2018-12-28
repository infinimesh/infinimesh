<template>
  <v-container>
    <v-layout
      row wrap
      >
      <v-card
        flat
      >
      <v-card-title
        primary-title
        class="body-2"
      >
        Device information
      </v-card-title>
      <v-card-text>
        Enabled: {{ device.enabled }}
      </v-card-text>
      <v-card-text>
        Id: {{ device.id }}
      </v-card-text>
      <v-card-text>
        <v-chip
         v-for="(tag, i) in device.tags"
         :key="i"
         small
        >
          {{ tag }}
        </v-chip>
      </v-card-text>
      </v-card>
    </v-layout>
  </v-container>
</template>

<script>
import { APIMixins } from "../mixins/APIMixins";

export default {
  mixins: [APIMixins],
  data() {
    return {
      device: {},
      checkbox: false,
      id: this.$route.params.id,
      headers: ["Active", "Id", "Name", "Location", "Tags"]
    };
  },
  mounted() {
    this.getRemoteDevice();
  }
};
</script>

<style lang="css">
</style>
