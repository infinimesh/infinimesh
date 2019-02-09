<template>
  <v-container>
    <v-card>
      <v-card-title>
        <h1 class="mb-3">Shadow - {{ device.id }}</h1>
      </v-card-title>
      <v-divider></v-divider>
      <v-layout row wrap>
        <v-flex>
          <v-card
            flat
            class="pb-3"
            min-width="400px"
          >
            <v-card-title>
              <h2>Reported state</h2>
            </v-card-title>
            <div
              id="scrollableCard"
            >
            <v-card>
            <v-card
              v-for="(response, index) in shadowMessages"
              :key="index"
            >
              <v-card-text>
                <strong>Timestamp</strong>: {{ response.result.reportedDelta.timestamp }}
                <v-spacer></v-spacer>
                <strong>Data</strong>: {{ response.result.reportedDelta.data }}
              </v-card-text>
            </v-card>
            <v-card-text>
              <strong>Initial timestamp</strong>: {{ shadow.initialState.timestamp }}
              <v-spacer></v-spacer>
              <strong>Initial data</strong>: {{ shadow.initialState.data }}
            </v-card-text>
          </v-card>
            </div>
          </v-card>
        </v-flex>
        <v-divider
          vertical
        ></v-divider>
        <component
          :is="activeComp"
          @edit="activeComp='Update'"
          @close="activeComp='DeviceInfo'"
        ></component>
      </v-layout>
    </v-card>
  </v-container>
</template>

<script>
import DeviceInfo from "../components/DeviceInfo.vue";
import Update from "../components/Update.vue";

export default {
  data() {
    return {
      device: {},
      shadow: {
        initialState: {
          data: "No data received",
          timestamp: "N/A"
        }
      },
      activeComp: DeviceInfo,
      id: this.$route.params.id,
      messages: []
    };
  },
  computed: {
    shadowMessages() {
      return this.$store.getters.getShadowMessages;
    }
  },
  methods: {
    getInitialShadow() {
      this.$store
        .dispatch("fetchInitialShadow", this.id)
        .then(() => {
          this.shadow.initialState = this.$store.getters.getInitialShadow;
        })
        .catch(e => console.log(e));
    },
    getDevice() {
      this.$store
        .dispatch("fetchDevices")
        .then(() => {
          this.device = this.$store.getters.getDevice(this.id);
          this.checkbox = this.$store.getters.getDevice(this.id).enabled;
        })
        .catch(e => console.log(e));
    }
  },
  created() {
    this.getDevice();
    this.getInitialShadow();
    this.$store.dispatch("connectToShadow", this.id);
  },
  components: {
    DeviceInfo,
    Update
  }
};
</script>

<style lang="css" scoped>
  #scrollableCard {
  height: 400px;
  overflow-y: auto;
  }
</style>
