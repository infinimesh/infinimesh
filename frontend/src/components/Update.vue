<template>
  <div>
    <v-card-title
      primary-title
    >
      <h2>Update device information</h2>
    </v-card-title>
    <v-card-text>
      <v-text-field
        v-model="tag"
        label="Device tags"
        clearable
        box
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
      <v-checkbox
       label="Device enabled"
       v-model="checkbox"
       class="mt-5"
      >
      </v-checkbox>
    </v-card-text>
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
    <v-card-actions>
      <v-layout
        row
        wrap
      >
        <div>
          <v-btn
            round
            class="mr-4 mb-3"
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
           class="mb-3"
           @click="updateDevice()"
         >
           Update device
         </v-btn>
         </div>
     </v-layout>
    </v-card-actions>
  </div>
</template>

<script>
export default {
  data() {
    return {
      device: {
        enabled: false,
        id: "",
        tags: []
      },
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
              this.$emit("close");
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
  created() {
    this.$store
      .dispatch("fetchDevices")
      .then(() => {
        this.device = this.$store.getters.getDevice(this.id);
        this.checkbox = this.$store.getters.getDevice(this.id).enabled;
      })
      .catch(e => console.log(e));
  }
};
</script>

<style lang="css">
</style>
