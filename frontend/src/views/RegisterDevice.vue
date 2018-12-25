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
             round
             color="secondary lighten-2"
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
import UploadButton from 'vuetify-upload-button';

export default {
  data() {
    return {
      id: "",
      tag: "",
      tags: [],
      certificate: {
        pem_data: `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURpRENDQW5DZ0F3SUJBZ0lKQU1OTk9LaE05ZXlPTUEwR0NTcUdTSWIzRFFFQkN3VUFNRmt4Q3pBSkJnTlYKQkFZVEFrRlZNUk13RVFZRFZRUUlEQXBUYjIxbExWTjBZWFJsTVNFd0h3WURWUVFLREJoSmJuUmxjbTVsZENCWAphV1JuYVhSeklGQjBlU0JNZEdReEVqQVFCZ05WQkFNTUNXeHZZMkZzYUc5emREQWVGdzB4T0RBNE1EWXlNVFU0Ck5UUmFGdzB5T0RBNE1ETXlNVFU0TlRSYU1Ga3hDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWwKTFZOMFlYUmxNU0V3SHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXhFakFRQmdOVgpCQU1NQ1d4dlkyRnNhRzl6ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTHEyCjVUMms5Ujk4aldtR1hqZUZyK2l1dGlndHV3STlUUTVDUTErMlJoOXNZcEV6eVpTZUhtMi9rZU1taGZ1TEQ5dnYKcU42a0hXV0FybXFMRkdaN01NMjh3cHNYT3hNZ0s1VUNsbVliOTVqWVVlbUtRbjZvcFNZQ25hcHZVajZVaHVCbwpjcGc3bTZlTHlzRzBXTVFaQW8xTEMyZU1JUUdUQ0JtWHVWRmFrUkwrMENGamFENWQ0K1ZKVUtodk1QTTV4cHR5CnFEMkJrOUtYTkhnUzh1WDhZeHhlMHRCK3A2UDYwS2d2OSt5V0NybTJSVVYvenVTbFhYNjluVUUvVnJlelNkR24KYy90VlNJY3NwaVhUcERsS2lITFBvWWZMODN4d01yd2c0WTFFVVREemtBa3U5OHVwc3MrR0RhbGtKYVNsZHk2NwpKSkxUczk0WmdHNXZKVFpQSmUwQ0F3RUFBYU5UTUZFd0hRWURWUjBPQkJZRUZKT0Vtb2I2cHRobkZacTJsWnpmCjM4d2ZRWmhwTUI4R0ExVWRJd1FZTUJhQUZKT0Vtb2I2cHRobkZacTJsWnpmMzh3ZlFaaHBNQThHQTFVZEV3RUIKL3dRRk1BTUJBZjh3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUpVaUFHSlFiSFBNZVlXaTRiT2hzdVVydkhoUAptTi9nNG53dGprQWl1NlE1UU9IeTF4VmRHelI3dTZyYkhaRk1tZElyVVBRLzVta3FKZFpuZGw1V1NoYnZhRy84CkkwVTNVcTBCM1h1ZjBmMVBjbjI1aW9UaitVN1BJVVlxV1FYdmpOMVlubHNVamNiUTdDUTJFT0hLbU5BN3YyZmcKT21XckJBcDRxcU9hRUtXcGcwTjlmWklDYjdnNGtsT05RT3J5QWFaWWNiZUNCd1h5ZzBiYUNaTFhmSnphdG40MQpYa3JyMG5Wd2VYaUVFazVCb3NOMjBGeUZaQmVrcGJ5MTF0aDJNMVhrc0FyTFRXUTQxSUwxVGZXS0pBTERaZ1BMCkFYOTlJS0VMelZUc25ka2ZGOG1MVldacjFPb2I3c29UVlhmT0kvVkJuMWUrM3FrVXJLOTRKWXRZajA0PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`,
        algorithm: "testalg"
      },
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
    addTag() {
      if (this.tag) {
        this.tags.push(this.tag);
        this.tag = "";
      }
    },
    setId() {
      this.id = "id-" + Math.random();
    },
    fileChanged(file) {
      let reader = new FileReader();
      let that = this;
      reader.onload = function(e) {
        that.certificate.pem_data = reader.result;
      }
      reader.readAsText(file);
     },
    register(enabled) {
      this.addTag();
      this.setId();
      this.$http
        .post("http://localhost:8081/devices", {
          id: this.id,
          enabled,
          certificate: this.certificate,
          tags: this.tags
        })
        .then((response) => {
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
  }
};
</script>

<style lang="css">
</style>
