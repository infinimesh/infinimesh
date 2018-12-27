<template>
  <v-container>
    <v-layout column wrap md9 lg6 xl4>
      <h1 class="mb-3">Update device information</h1>
      <v-flex>
        <v-card
          class="pa-3"
        >
          <v-card-title
            primary-title
            class="body-2"
          >
            Device ID
          </v-card-title>
          <v-card-text>
            {{ device.id }}
          </v-card-text>
        </v-card>
        <v-card
          class="mt-2 pa-3"
        >
          <v-checkbox
            label="Device enabled"
            v-model="checkbox"
          ></v-checkbox>
        </v-card>
        <v-card
        class="mt-2 pa-3"
        >
          <v-text-field
            v-model="tag"
            label="Device tags"
            clearable
            v-on:keyup.enter="addTag($event)"
          >
          </v-text-field>
          <v-chip
           v-for="(tag, i) in device.tags"
           :key="i"
           small
          >
             {{ tag }}
            <v-icon
              class="ml-1"
              small
              @click="device.tags.splice(i, 1)"
              style="color: grey"
            >
              cancel
            </v-icon>
         </v-chip>
        </v-card>
        <v-alert
          :value="messageSuccess.value"
          type="success"
          icon="check_circle"
        >
          {{ messageSuccess.message }}
        </v-alert>
        <v-alert
         :value="messageFailure.value"
         type="error"
         icon="error"
        >
          {{ messageFailure.value }}: {{ messageFailure.error }}
        </v-alert>
        <v-layout
          row
          wrap
          class="mt-3"
        >
          <div>
            <v-btn
              round
              class="mr-5"
              to="/devices"
            >
              <v-icon
                left
              >
                chevron_left
              </v-icon>
              Return
            </v-btn>
          </div>
          <div>
           <v-btn
             round
             color="primary"
             @click="updateDevice()"
           >
             Update device
           </v-btn>
           </div>
           </div>
       </v-layout>
      </v-flex>
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
      headers: ["Active", "Id", "Name", "Location", "Tags"],
      tag: "",
      messageSuccess: {
        message: "Update successful",
        value: false
      },
      messageFailure: {
        message: "Error in updating device",
        value: false,
        error: ""
      }
    };
  },
  methods: {
    addTag() {
      if (this.tag) {
        this.device.tags.push(this.tag);
        this.tag = "";
      }
    },
    updateDevice() {
      this.addTag();
      this.$http
        .patch("http://localhost:8081/devices/" + this.id, {
          enabled: this.checkbox,
          tags: this.device.tags
        })
        .then(response => {
          if (response.status === 200) {
            this.messageSuccess.value = true;
            setTimeout(() => {
              this.messageSuccess.value = false;
              this.$router.push("/devices");
            }, 1000);
          }
        })
        .catch(e => {
          this.messageFailure.value = true;
          this.messageFailure.error = e;
        });
      }
    },
    resetForm() {
      this.id = "";
      this.device.tags = [];
      this.enabled = false;
    },
  mounted() {
    this.getRemoteDevice();
  }
};
</script>

<style lang="css">
</style>
