<template>
  <v-container>
    <v-layout column wrap md9 lg6 xl4>
      <h1 class="mb-3">Register new device</h1>
      <v-flex>
        <v-card
          class="pa-3"
        >
          <v-text-field
            label="Device Id"
            v-model="id"
            :rules="idRules"
          ></v-text-field>
        </v-card>
        <v-card
          class="mt-2 pa-3"
        >
          <v-checkbox
            label="Device enabled"
            v-model="checkbox"
            clearable
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
           v-for="(tag, i) in tags"
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
        </v-card>
        <v-card
        class="mt-2 pa-3"
        >
          <v-layout row wrap>
            <v-flex>
              <v-textarea
               v-model="certificate.pem_data"
               auto-grow
               clearable
               label="Certificate"
               rows="1"
               >
              </v-textarea>
            </v-flex>
            <v-flex
             class="ml-3"
            >
             <upload-button
               round
               color="secondary lighten-2"
               class="white--text"
               :fileChangedCallback="fileChanged"
             >
               <template slot="icon">
                 <v-icon
                   class="ml-2"
                   style="color: white"
                 >
                   cloud_upload
                 </v-icon>
               </template>
             </upload-button>
            </v-flex>
          </v-layout>
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
             round
             color="primary"
             dark
             @click="register()"
           >
             Register device
           </v-btn>
         </div>
       </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import UploadButton from "vuetify-upload-button";
import { storeRemoteDevices } from "../mixins/APIMixins";

export default {
  mixins: [storeRemoteDevices],
  data() {
    return {
      checkbox: true,
      id: "",
      idRules: [
        v => !!v || "Id is required",
        v => !v.match(/\s/) || "No whitespace allowed",
        v =>
          !this.$store.getters.getDevice(v) || "This device Id already exists"
      ],
      tag: "",
      tags: [],
      certificate: {
        pem_data: `-----BEGIN CERTIFICATE-----\nMIIDiDCCAnCgAwIBAgIJAMNNOKhM9eyOMA0GCSqGSIb3DQEBCwUAMFkxCzAJBgNV\nBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX\naWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xODA4MDYyMTU4\nNTRaFw0yODA4MDMyMTU4NTRaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21l\nLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNV\nBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALq2\n5T2k9R98jWmGXjeFr+iutigtuwI9TQ5CQ1+2Rh9sYpEzyZSeHm2/keMmhfuLD9vv\nqN6kHWWArmqLFGZ7MM28wpsXOxMgK5UClmYb95jYUemKQn6opSYCnapvUj6UhuBo\ncpg7m6eLysG0WMQZAo1LC2eMIQGTCBmXuVFakRL+0CFjaD5d4+VJUKhvMPM5xpty\nqD2Bk9KXNHgS8uX8Yxxe0tB+p6P60Kgv9+yWCrm2RUV/zuSlXX69nUE/VrezSdGn\nc/tVSIcspiXTpDlKiHLPoYfL83xwMrwg4Y1EUTDzkAku98upss+GDalkJaSldy67\nJJLTs94ZgG5vJTZPJe0CAwEAAaNTMFEwHQYDVR0OBBYEFJOEmob6pthnFZq2lZzf\n38wfQZhpMB8GA1UdIwQYMBaAFJOEmob6pthnFZq2lZzf38wfQZhpMA8GA1UdEwEB\n/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAJUiAGJQbHPMeYWi4bOhsuUrvHhP\nmN/g4nwtjkAiu6Q5QOHy1xVdGzR7u6rbHZFMmdIrUPQ/5mkqJdZndl5WShbvaG/8\nI0U3Uq0B3Xuf0f1Pcn25ioTj+U7PIUYqWQXvjN1YnlsUjcbQ7CQ2EOHKmNA7v2fg\nOmWrBAp4qqOaEKWpg0N9fZICb7g4klONQOryAaZYcbeCBwXyg0baCZLXfJzatn41\nXkrr0nVweXiEEk5BosN20FyFZBekpby11th2M1XksArLTWQ41IL1TfWKJALDZgPL\nAX99IKELzVTsndkfF8mLVWZr1Oob7soTVXfOI/VBn1e+3qkUrK94JYtYj04=\n-----END CERTIFICATE-----`,
        algorithm: "testalg"
      },
      messageSuccess: {
        message: "Your device has been registered",
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
    addTag() {
      if (this.tag) {
        this.tags.push(this.tag);
        this.tag = "";
      }
    },
    fileChanged(file) {
      let reader = new FileReader();
      let that = this;
      reader.onload = function() {
        that.certificate.pem_data = reader.result;
      };
      reader.readAsText(file);
    },
    register() {
      this.addTag();
      this.$http
        .post("devices", {
          id: this.id,
          enabled: this.checkbox,
          certificate: this.certificate,
          tags: this.tags
        })
        .then(response => {
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
      this.tags = [];
      this.enabled = false;
    }
  },
  components: {
    UploadButton
  },
  beforeMount() {
    this.storeRemoteDevices();
  }
};
</script>

<style lang="css">
</style>
