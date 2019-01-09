<template>
  <v-container>
    <v-card
      flat
    >
    <v-layout
      row
      wrap
    >
      <v-flex>
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
      </v-flex>
      <v-flex
        offset-lg7
      >
        <v-icon
          style="cursor: pointer"
          @click="$emit('edit')"
        >
          edit
        </v-icon>
      </v-flex>
    </v-layout>
    </v-card>
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
