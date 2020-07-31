<template>
  <a-row type="flex" justify="space-between">
    <a-col :span="24">
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

export default Vue.component("device-actions", {});
</script>
