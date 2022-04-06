<template>
  <n-space align="center" justify="center" class="fullscreen">
    <LoginCard />
  </n-space>
</template>

<script setup>
import { NSpace } from "naive-ui";
import LoginCard from "@/components/login-card.vue";
</script>

<style>
.fullscreen {
  height: 99vh !important;
  width: 99vw !important;
}
</style>