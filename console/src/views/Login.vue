<template>
  <n-space align="center" justify="center" class="fullscreen" :wrap-item="false">
    <LoginCard />
  </n-space>
</template>

<script setup>
import { NSpace } from "naive-ui";
import LoginCard from "@/components/login-card.vue";
</script>

<style>
.fullscreen {
  height: 100vh !important;
  width: 100vw !important;
  padding: 15px;
  box-sizing: border-box;
}
</style>