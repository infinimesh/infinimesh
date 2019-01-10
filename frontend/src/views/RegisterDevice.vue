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
        pem_data: `-----BEGIN CERTIFICATE-----\nMIIDxzCCAq+gAwIBAgIUSQ5nH2Cgw2Q3lKyQau12bOGzV8YwDQYJKoZIhvcNAQEL\nBQAwczELMAkGA1UEBhMCREUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM\nGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEKMAgGA1UEAwwBKjEgMB4GCSqGSIb3\nDQEJARYRam9lQGluZmluaW1lc2guaW8wHhcNMTkwMTEwMjAwMTI1WhcNMjkwMTA3\nMjAwMTI1WjBzMQswCQYDVQQGEwJERTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8G\nA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMQowCAYDVQQDDAEqMSAwHgYJ\nKoZIhvcNAQkBFhFqb2VAaW5maW5pbWVzaC5pbzCCASIwDQYJKoZIhvcNAQEBBQAD\nggEPADCCAQoCggEBALq25T2k9R98jWmGXjeFr+iutigtuwI9TQ5CQ1+2Rh9sYpEz\nyZSeHm2/keMmhfuLD9vvqN6kHWWArmqLFGZ7MM28wpsXOxMgK5UClmYb95jYUemK\nQn6opSYCnapvUj6UhuBocpg7m6eLysG0WMQZAo1LC2eMIQGTCBmXuVFakRL+0CFj\naD5d4+VJUKhvMPM5xptyqD2Bk9KXNHgS8uX8Yxxe0tB+p6P60Kgv9+yWCrm2RUV/\nzuSlXX69nUE/VrezSdGnc/tVSIcspiXTpDlKiHLPoYfL83xwMrwg4Y1EUTDzkAku\n98upss+GDalkJaSldy67JJLTs94ZgG5vJTZPJe0CAwEAAaNTMFEwHQYDVR0OBBYE\nFJOEmob6pthnFZq2lZzf38wfQZhpMB8GA1UdIwQYMBaAFJOEmob6pthnFZq2lZzf\n38wfQZhpMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAKT6Eien\nI/ngNjd5ZPW1JCLPoWL8nYO4MRICZgbi1d/Kp0478jJWzYYSffH4qustIZFKzzgB\n5Qa0JTKWdpS4dyo1jgaMfuExOvFkdtLtztDNkILBHwgSOeN4q1s2343AMpRqPulB\nhtU0vDyYIvDwRZly/anMgAeNMH2vFee6Z46w1BrO24LWafV3cBMEKTPepHEstpLn\n4t+tODI3XIDLO3Lj8Lcm8FVT+m6pW3iEcRGikRNLjD1UPDZvHHXlO3dGrsuuDCca\np/2Jrz3RlDBNbwOAjzSWGcXfKAhEtp3eXyE7gvegKC1r2iB0Nu/GE/I3B7oxl4hp\npoX4zHWffhe2X+Q=\n-----END CERTIFICATE-----`,
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
