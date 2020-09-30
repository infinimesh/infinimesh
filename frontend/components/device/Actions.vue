<template>
  <a-row type="flex" :gutter="10">
    <a-col>
      <a-popconfirm
        :title="
          `Are you sure ${
            device.enabled ? 'disabling' : 'enabling'
          } this device?`
        "
        ok-text="Yes"
        cancel-text="No"
        @confirm="$emit('toogle')"
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
  </a-row>
</template>

<script>
import Vue from "vue";

export default Vue.component("device-actions", {
  props: {
    device: {
      required: true
    }
  }
});
</script>

<style scoped>
.anticon {
  color: var(--icon-color-dark);
}
</style>
