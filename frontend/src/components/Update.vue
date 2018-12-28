<template>
  <v-container>
    <v-layout column wrap md9 lg6 xl4>
      <v-card
        flat
      >
        <v-card-title
          primary-title
          class="body-2"
        >
          Update device information
        </v-card-title>
        <v-flex>
          <v-card
            flat
            class="ml-3"
          >
            <v-checkbox
              label="Device enabled"
              v-model="checkbox"
            ></v-checkbox>
          </v-card>
          <v-card
            class="pb-3 ml-3"
            flat
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
          <v-card>
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
          </v-card>          
          <v-layout
            row
            wrap
            class="mt-3 ml-2"
          >
            <div>
              <v-btn
                round
                class="mr-4"
                @click="$emit('close')"
                small
              >
                <v-icon
                  left
                >
                  close
                </v-icon>
                Abort
              </v-btn>
            </div>
            <div>
             <v-btn
               round
               small
               @click="updateDevice()"
             >
               Update device
             </v-btn>
             </div>
         </v-layout>
        </v-flex>
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
        .patch("devices/" + this.id, {
          enabled: this.checkbox,
          tags: this.device.tags
        })
        .then(response => {
          if (response.status === 200) {
            this.$store.dispatch("updateDevice", {
              id: this.id,
              enabled: this.checkbox,
              tags: this.device.tags
            });
            this.messageSuccess.value = true;
            setTimeout(() => {
              this.messageSuccess.value = false;
              this.$emit('close');
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
