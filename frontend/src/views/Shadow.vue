<template>
  <v-container>
    <h1 class="mb-3">Device overview</h1>
    <v-card>
      <v-layout row wrap>
        <v-flex>
          <v-card
            flat
          >
            <v-card-title
            primary-title
            class="body-2"
            >
              Device shadow
            </v-card-title>
          </v-card>
        </v-flex>
        <v-divider
          vertical
        ></v-divider>
        <v-flex>
          <v-card
            flat
          >
            <component
              :is="activeComp"
              @edit="activeComp='Update'"
              @close="activeComp='DeviceInfo'"
            ></component>
          </v-card>
        </v-flex>
      </v-layout>
    </v-card>
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
      activeComp: DeviceInfo,
      id: this.$route.params.id
    };
  },
  methods: {
    connectToShadow(id) {
      let xhr = new XMLHttpRequest();

      xhr.open("GET", `http://localhost:8081/devices/${id}/shadow/reported`, true)
      xhr.onprogress = function () {
       console.log("PROGRESS:", xhr.responseText)
      }
      xhr.send()
    }
  },
  mounted() {
    this.getRemoteDevice();
    this.connectToShadow(this.id);
  },
  components: {
    DeviceInfo,
    Update
  }
};
</script>

<style lang="css">
</style>
