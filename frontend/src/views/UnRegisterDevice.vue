<template>
  <v-container xs12>
    <v-layout column wrap>
      <h1 class="mb-3">Unregister a device</h1>
      <v-flex>
        <v-card>
          <v-card-title primary-title>
            Are you sure you want to unregister device with Id: {{ id }}
          </v-card-title>
          <v-card-text>
            This will prevent clients from further communication with this
            device.
          </v-card-text>
          <v-alert
            :value="messageSuccess.value"
            type="success"
            icon="check_circle"
          >
            {{ messageSuccess.message }}
          </v-alert>
          <v-alert :value="messageFailure.value" type="error" icon="error">
            {{ messageFailure.value }}: {{ messageFailure.error }}
          </v-alert>
          <v-card-actions>
            <v-layout row wrap>
              <div>
                <v-btn round class="mr-5" to="/devices">
                  <v-icon left>
                    chevron_left
                  </v-icon>
                  Return
                </v-btn>
              </div>
              <div>
                <v-btn round color="primary" dark @click="unRegisterDevice(id)">
                  Unregister device
                </v-btn>
              </div>
            </v-layout>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      id: this.$route.params.id,
      messageSuccess: {
        message: "Your device has been unregistered",
        value: false
      },
      messageFailure: {
        message: "Error in unregistering device",
        value: false,
        error: ""
      }
    };
  },
  computed: {
    ...mapGetters({namespace: "getNamespace"})
  },
  methods: {
    unRegisterDevice(id) {
      this.$store.dispatch("unRegisterDevice", id);

      this.$http
        .delete("devices/" + id)
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
  }
};
</script>

<style lang="css"></style>
