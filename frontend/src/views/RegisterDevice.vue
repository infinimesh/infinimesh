<template>
  <v-container>
    <v-layout column wrap md9 lg6 xl4>
      <h1 class="mb-3">Register new device</h1>
      <v-flex>
        <v-text-field
          v-model="tag"
          label="Device tags"
          clearable
          v-on:keyup.enter="addTag($event)"
        >
        </v-text-field>
        <v-chip
         v-for="(tag, key, i) in tags"
         :key="i"
         small
        >
           {{ tag }}
          <v-icon
            class="ml-1"
            small
            @click="tags.splice(i, 1)"
            style="color: grey"
          >
            cancel
          </v-icon>
       </v-chip>
        <v-textarea
         v-model="certificate"
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
             @click="register(true)"
           >
             Register and activate</v-btn>
           </div>
           <div>
           <v-btn
             round color="secondary lighten-2"
             dark
             @click="register(false)"
           >
             Register and don't activate
          </v-btn>
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
      tag: "",
      tags: [],
      certificate: "",
      messageSuccess: {
        message: "Your device has been enabled",
        value: false
      },
      messageFailure: {
        message: "Error in registering device",
        value: false,
        error: ""
      }
    };
  },
  methods: {
    addTag(event) {
      this.tags.push(event.target.value);
      this.tag = "";
    },
    register(enabled) {
      const id = this.tags[0] + Math.random();
      let newDevice = {};
      newDevice = {
        id,
        enabled,
        tags: this.tags,
        certificate: this.certificate
      };
      this.addRemote(newDevice);
    },
    addRemote(device) {
      this.$http
        .post("http://localhost:8081/devices", device)
        .then((response) => {
          console.log(response);
          this.$store.dispatch("addDevice", newDevice);
          this.resetForm();
          this.messageSuccess.value = true;
          setTimeout(() => (this.messageSuccess.value = false), 5000);
        })
        .catch(e => {
          this.messageFailure.value = true;
          this.messageFailure.error = e;
        });
    },
    resetForm() {
      this.name = "";
      this.description = "";
      this.location = "";
      this.tags = [];
      this.certificate = "";
      this.enabled = false;
    }
  }
};
</script>

<style lang="css">
</style>
