<template>
  <n-space justify="space-between" :style="{padding: '15px'}" align="center">
    <span class="infinimesh">infinimesh<span class="copyright">©</span></span>
    <user-details/>
  </n-space>
</template>

<script setup>
import { NSpace } from "naive-ui"
import UserDetails from "@/components/dashboard/user.vue"
</script>

<style>
span.infinimesh {
  font-size: 2vh;
  font-family: 'Exo', sans-serif;
}
span.copyright {
  vertical-align: super;
  font-size: 0.8vh;
}
</style>