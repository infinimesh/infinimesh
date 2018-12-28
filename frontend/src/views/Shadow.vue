<template>
  <v-container>
    <h1 class="mb-3">Device overview</h1>
    <v-layout row wrap>
      <v-flex>
        <v-card>
          <v-card-title
          primary-title
          class="body-2"
          >
            Device shadow
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex>
        <v-card
        <v-icon
        >
          edit
        </v-icon>
        >
          <component
            :is="activeComp"
          ></component>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { APIMixins } from "../mixins/APIMixins";
import DeviceInfo from "../components/DeviceInfo.vue";
import Update from "../components/Update.vue";

export default {
  mixins: [APIMixins],
  data() {
    return {
      activeComp: DeviceInfo
    }
  },
  mounted() {
    this.getRemoteDevice();
  },
  components: {
    DeviceInfo,
    Update
  }
}
</script>

<style lang="css">
</style>
