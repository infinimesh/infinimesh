<template>
  <div>
    <v-card
      flat
      max-width="500"
    >
      <v-layout row wrap>
        <v-card-title
          primary-title
        >
          <h2>Device information</h2>
        </v-card-title>
        <v-layout
          align-end
          justify-end
        >
          <v-icon
            style="cursor: pointer"
            @click="$emit('edit')"
            class="ma-3"
          >
            edit
          </v-icon>
        </v-layout>
          <v-card-text
          >
            <v-icon
              v-if="device.enabled"
              color="green"
              class="mr-2"
            >
              check_circle
            </v-icon>
            <v-icon
              v-else
              color="grey"
              class="mr-2"
            >
              block
            </v-icon>
            {{ (device.enabled) ? "Device enabled" : "Device disabled" }}
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
      </v-layout>
    </v-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      device: {
        enabled: false,
        id: "",
        tags: []
      },
      checkbox: false,
      id: this.$route.params.id,
      headers: ["Active", "Id", "Name", "Location", "Tags"]
    };
  },
  created() {
    this.$store
      .dispatch("fetchDevices")
      .then(() => {
        this.device = this.$store.getters.getDevice(this.id);
        this.checkbox = this.$store.getters.getDevice(this.id).enabled;
      })
      .catch(e => console.log(e));
  }
};
</script>

<style lang="css">
</style>
