<template>
  <a-row type="flex" justify="space-between">
    <a-col :span="24">
      <a-button type="danger" @click="$emit('delete')">Delete</a-button>
    </a-col>
  </a-row>
</template>

<script>
import Vue from "vue";

export default Vue.component("device-actions", {});
</script>
