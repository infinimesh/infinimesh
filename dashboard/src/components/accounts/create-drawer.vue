<template>
  <n-button @click="show = true" type="success" dashed>
    <template #icon>
      <n-icon><add-outline /></n-icon>
    </template>
    Create Account
  </n-button>
</template>

<script setup>
import { NButton, NIcon } from "naive-ui";
import { AddOutline } from "@vicons/ionicons5";
</script>