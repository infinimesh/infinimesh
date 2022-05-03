<template>
  <n-space justify="space-between" :style="{ padding: '15px' }" align="center">
    <span class="infinimesh">infinimesh<span class="copyright">©</span></span>
    <ns-selector />
    <theme-picker style="margin-left: 5px" />
    <user-details />
  </n-space>
</template>

<script setup>
import { NSpace } from "naive-ui";
import UserDetails from "@/components/dashboard/user.vue";
import NsSelector from "@/components/core/ns-selector.vue";
import ThemePicker from "@/components/core/theme-picker.vue";
</script>

<style>
span.infinimesh {
  font-size: 2vh;
  font-family: "Exo 2", sans-serif;
}

span.copyright {
  vertical-align: super;
  font-size: 0.8vh;
}
</style>