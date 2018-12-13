<template>
  <v-container xs9>
    <v-layout column wrap>
      <h1 class="mb-3">Register new device</h1>
      <v-flex>
        <v-text-field
          v-model="name"
          label="Device name"
          outline
          clearable
        ></v-text-field>
        <v-text-field
          v-model="description"
          label="Device description"
          outline
          clearable
        ></v-text-field>
        <v-text-field
          v-model="location"
          label="Device location"
          outline
          clearable
        ></v-text-field>
        <v-text-field
          v-model="tags"
          label="Device tags"
          outline
          clearable
        ></v-text-field>
        <v-textarea
           v-model="certificate"
           auto-grow
           outline
           label="Certificate"
           rows="1"
         ></v-textarea>
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
           <div >
           <v-btn
           round color="primary"
           dark
           @click="register(true)"
           >
           Register and activate</v-btn>
           </div>
           <div >
           <v-btn
           round color="secondary lighten-2"
           dark
           @click="register(false)"
           >
           Register and don't activate</v-btn>
           </div>
           <p v-for="device in devices">{{ device }}</p>
         </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      name: "",
      description: "",
      location: "",
      tags: "",
      certificate: "",
      activated: false,
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
  computed: {
    devices() {
      return this.$store.getters.getAllDevices;
    }
  },
  methods: {
    register(activate) {
    const deviceId = this.name + Math.random();

    this.activated = activate;

    let newDevice = {};

    newDevice[`${deviceId}`] = {
      status: this.activated,
      name: this.name,
      description: this.description,
      location: this.location,
      tags: this.tags,
      certificate: this.certificate
    };

    this.$store.dispatch("addDevice", newDevice);

    this.resetForm();
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
    },
    resetForm() {
      this.name = "";
      this.description = "";
      this.location = "";
      this.tags = "";
      this.certificate = "";
      this.activated = false;
    }
  }
};
</script>

<style lang="css">
</style>
