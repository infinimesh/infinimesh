<template>
  <v-container>
    <h1 class="mb-3">{{ device.id }} - Device overview</h1>
    <v-card>
      <v-layout row wrap>
        <v-flex>
          <v-card
            flat
            class="pb-3"
          >
            <v-card-title
            primary-title
            >
              <h2>Reported state</h2>
            </v-card-title>
            <div
              id="scrollableCard"
            >
            <v-card>
            <v-card
              v-for="(response, index) in messages"
              :key="index"
            >
              <v-card-text>
                <strong>Timestamp</strong>: {{ response.result.reportedDelta.timestamp }}
                <v-spacer></v-spacer>
                <strong>Data</strong>: {{ response.result.reportedDelta.data }}
              </v-card-text>
            </v-card>
            <v-card-text>
              <strong>Initial timestamp</strong>: {{ shadow.initialTimestamp }}
              <v-spacer></v-spacer>
              <strong>Initial data</strong>: {{ shadow.initialState }}
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
import { APIMixins } from "../mixins/APIMixins";
import DeviceInfo from "../components/DeviceInfo.vue";
import Update from "../components/Update.vue";
import Vue from "vue";

export default {
  mixins: [APIMixins],
  data() {
    return {
      device: {},
      shadow: {
        initialState: "No data received",
        initialTimestamp: "N/A"
      },
      activeComp: DeviceInfo,
      id: this.$route.params.id,
      messages: []
    };
  },
  methods: {
    connectToShadow(id) {
      let xhr = new XMLHttpRequest();
      let that = this;

      setTimeout(() => {
        xhr.open(
          "GET",
          Vue.http.options.root + `/devices/${id}/shadow/reported`,
          true
        );
        xhr.onprogress = function() {
          let jsonObjects = [];
          let obj = "";
          that.messages = [];

          jsonObjects = xhr.responseText.replace(/\n$/, "").split(/\n/);
          for (obj of jsonObjects) {
            that.messages.push(JSON.parse(obj));
          }
          that.messages.reverse();
        };
        xhr.send();
      }, 1000);
    },
    getInitialShadow(id) {
      this.$http
        .get(`devices/${id}/shadow`)
        .then(response => {
          this.shadow.initialTimestamp = response.body.shadow.reported.timestamp;
          this.shadow.initialState = response.body.shadow.reported.data;
        })
        .catch(e => {
          console.log(e);
        });
    }
  },
  mounted() {
    this.getRemoteDevice();
    this.getInitialShadow(this.id);
    this.connectToShadow(this.id);
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
