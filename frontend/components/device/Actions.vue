<template>
  <a-row type="flex" :gutter="10">
    <a-col>
      <a-popconfirm
        :title="`Are you sure ${
          device.enabled ? 'disabling' : 'enabling'
        } this device?`"
        ok-text="Yes"
        cancel-text="No"
        @confirm="$emit('toggle')"
      >
        <a-button
          :type="device.enabled ? 'danger' : 'success'"
          icon="switcher"
          >{{ device.enabled ? "Disable" : "Enable" }}</a-button
        >
      </a-popconfirm>
    </a-col>
    <a-col>
      <a-popconfirm
        title="Are you sure deleting this device?"
        ok-text="Yes"
        cancel-text="No"
        @confirm="$emit('delete')"
      >
        <a-button type="danger" icon="delete">Delete</a-button>
      </a-popconfirm>
    </a-col>
    <a-col>
      <a-button
        :type="device.basic_enabled ? 'success' : 'danger'"
        icon="switcher"
        @click="() => (basic_enabled_visible = true)"
        >MQTT Basic Auth</a-button
      >
      <ToggleMQTTBasicAuth
        @toggle="$emit('toggle-basic')"
        @close="() => (basic_enabled_visible = false)"
        :device="device"
        :visible="basic_enabled_visible"
      />
    </a-col>
  </a-row>
</template>

<script>
import Vue from "vue";
import ToggleMQTTBasicAuth from "./ToggleMQTTBasicAuth";

export default Vue.component("device-actions", {
  components: {
    ToggleMQTTBasicAuth,
  },
  props: {
    device: {
      required: true,
    },
  },
  data() {
    return {
      basic_enabled_visible: false,
    };
  },
});
</script>

<style scoped>
.anticon {
  color: var(--icon-color-dark);
}
</style>
