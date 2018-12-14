<template>
  <v-container xs12>
    <v-layout column wrap>
      <h1 class="mb-3">Unregister device</h1>
      <v-flex>
        <v-card>
          <v-card-title primary-title>
            Are you sure you want to unregister device with Id: {{ deviceId }}
          </v-card-title>
          <v-card-text>
            This will prevent clients from further communication with this device.
          </v-card-text>
          <v-card-actions>
            <v-btn
            round color="primary"
            dark
            @click="unRegisterDevice(deviceId)"
            >
            Unregister device
            </v-btn>
          </v-card-actions>
        </v-card>
        <div>
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
        </div>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      deviceId: this.$route.params.id,
      messageSuccess: {
        message: "Your device has been activated",
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
    unRegisterDevice(deviceId) {

      this.$store.dispatch("unRegisterDevice", deviceId);

      this.$router.push("/devices");
      // this.$http
      //   .post("testdata.json", this.device)
      //   .then(() => {
      //     this.resetForm();
      //     this.messageSuccess.value = true;
      //     setTimeout(() => (this.messageSuccess.value = false), 5000);
      //   })
      //   .catch(e => {
      //     this.messageFailure.value = true;
      //     this.messageFailure.error = e;
      //   });
    }
  }
};
</script>

<style lang="css">
</style>
