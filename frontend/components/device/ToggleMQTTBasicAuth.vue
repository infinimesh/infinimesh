<template>
  <a-modal
    :visible="visible"
    title="Toggle MQTT Basic Auth"
    @ok="$emit('close')"
    @cancel="$emit('close')"
  >
    <a-row>
      <h3>MQTT Basic Auth Credentials</h3>
    </a-row>
    <a-row>
      <a-col :span="6">
        <h4>Username</h4>
      </a-col>
      <a-input :value="device.name" disabled>
        <a-icon
          slot="addonAfter"
          type="copy"
          @click="copyTextToClipboard(device.name)"
        />
      </a-input>
    </a-row>
    <a-row>
      <a-col :span="6">
        <h4>Password</h4>
      </a-col>
      <a-input :value="fingerprint" disabled>
        <a-icon
          slot="addonAfter"
          type="copy"
          @click="copyTextToClipboard(fingerprint)"
        />
      </a-input>
    </a-row>
    <a-row type="flex" justify="space-around" style="margin-top: 15px">
      <a-col>
        <a-popconfirm
          :title="`Are you sure ${
            device.basic_enabled ? 'disabling' : 'enabling'
          } MQTT Basic Auth for this device?`"
          ok-text="Yes"
          cancel-text="No"
          @confirm="$emit('toggle')"
        >
          <a-button
            :type="device.basic_enabled ? 'success' : 'danger'"
            icon="switcher"
            >{{ device.basic_enabled ? "Disable" : "Enable" }}</a-button
          >
        </a-popconfirm>
      </a-col>
    </a-row>
  </a-modal>
</template>

<script>
import Vue from "vue";
import Clipboard from "@/mixins/clipboard";

export default Vue.component("toggle-mqtt-basic-auth", {
  mixins: [Clipboard],
  props: {
    visible: {
      required: true,
      type: Boolean,
    },
    device: {
      required: true,
    },
  },
  computed: {
    fingerprint() {
      console.log(this.device);
      if (!this.device.certificate) return "";
      let binary_string = atob(this.device.certificate.fingerprint);
      let len = binary_string.length;
      let hash = "";
      for (var i = 0; i < len; i++) {
        hash += binary_string
          .charCodeAt(i)
          .toString(16)
          .padStart(2, 0)
          .toUpperCase();
      }
      return hash;
    },
  },
});
</script>

<style scoped>
.ant-input-disabled {
  color: rgba(0, 0, 0, 0.9);
}
</style>