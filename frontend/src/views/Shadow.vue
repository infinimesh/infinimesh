<template>
  <div>
    <v-card-title>
      <h1 class="mb-3">{{ device.id }}</h1>
    </v-card-title>
    <v-divider></v-divider>
    <v-layout row wrap>
      <v-flex>
        <v-card flat class="pb-3" width="50%">
          <v-card-title>
            <h2>Reported state</h2>
          </v-card-title>
          <div id="scrollableCard">
            <v-card flat>
              <v-card
                v-for="(response, index) in shadowMessages"
                :key="index"
                flat
              >
                <v-card-text>
                  <strong>Timestamp</strong>:
                  {{ response.result.reportedDelta.timestamp }}
                  <v-spacer></v-spacer>
                  <strong>Data</strong>:
                  {{ response.result.reportedDelta.data }}
                </v-card-text>
              </v-card>
              <v-card-text>
                <strong>Initial timestamp</strong>:
                {{ shadow.initialState.timestamp }}
                <v-spacer></v-spacer>
                <strong>Initial data</strong>: {{ shadow.initialState.data }}
              </v-card-text>
            </v-card>
          </div>
        </v-card>
      </v-flex>
      <v-divider vertical></v-divider>
      <v-card flat width="50%">
        <v-layout row>
          <v-flex>
            <component :is="activeComp"></component>
          </v-flex>
          <v-flex max-width="30px">
            <v-toolbar card color="white">
              <v-spacer />
              <v-btn fab small @click="isEditing = !isEditing">
                <v-icon v-if="isEditing">mdi-close</v-icon>
                <v-icon v-else>mdi-pencil</v-icon>
              </v-btn>
            </v-toolbar>
          </v-flex>
        </v-layout>
      </v-card>
    </v-layout>
  </div>
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
      isEditing: false,
      id: this.$route.params.id,
      messages: []
    };
  },
  computed: {
    shadowMessages() {
      return this.$store.getters.getShadowMessages;
    },
    activeComp() {
      return this.isEditing ? "Update" : "DeviceInfo";
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
