<template>
  <v-container>
    <v-layout column wrap md9 lg6 xl4>
      <h1 class="mb-3">Update device information</h1>
      <v-flex>
        <v-text-field
          v-model="tag"
          label="Device tags"
          clearable
          v-on:keyup.enter="addTag($event)"
        >
        </v-text-field>
        <v-chip
         v-for="(tag, key, i) in device.tags"
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
        <v-textarea
         v-model="device.certificate"
         auto-grow
         clearable
         label="Certificate"
         rows="1"
         >
        </v-textarea>
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
        <v-layout row wrap>
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
             round color="primary"
             dark
             @click="update(true)"
           >
             Update</v-btn>
           </div>
       </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      id: this.$route.params.id,
      tag: "",
      device: {},
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
  beforeMount() {
    this.device = this.$store.getters.getDevice(this.id);
  },
  methods: {
    addTag() {
      if (this.tag) {
        this.device.tags.push(this.tag);
        this.tag = "";
      }
    },
    update(enabled) {
      this.addTag();
      this.$http
        .post("http://localhost:8081/devices", {
          id: this.id,
          enabled,
          certificate: this.device.certificate,
          tags: this.device.tags
        })
        .then((response) => {
          console.log(response);
          if (response.status === 200) {
            this.resetForm();
            this.messageSuccess.value = true;
            setTimeout(() => (this.messageSuccess.value = false), 5000);
          }
        })
        .catch(e => {
          this.messageFailure.value = true;
          this.messageFailure.error = e;
        });
    },
    resetForm() {
      this.id = "";
      this.device.tags = [];
      this.enabled = false;
    }
  }
};
</script>

<style lang="css">
</style>
