<template>
  <n-config-provider :theme="theme.it" :theme-overrides="theme.overrides" :hljs="hljs">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-global-style />
        <router-view />
        <n-watermark v-if="watermark" content="development preview" cross fullscreen :font-size="16" :line-height="16"
          :width="250" :height="150" :x-offset="12" :y-offset="80" :rotate="-15" />
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup>
import { ref, computed, inject } from "vue";
import {
  NConfigProvider,
  NGlobalStyle,
  NWatermark,
  NLoadingBarProvider,
  NMessageProvider,
  darkTheme,
  lightTheme,
} from "naive-ui";

import { storeToRefs } from 'pinia';
import { useAppStore } from "@/store/app"

import lightThemeOverrides from "@/assets/light-theme-overrides.json"
import darkThemeOverrides from "@/assets/dark-theme-overrides.json"

import hljs from "@/utils/hljs";

const store = useAppStore()
const { theme: pick, base_url } = storeToRefs(store)
const theme = computed(() => {
  return {
    it: pick.value === "dark" ? darkTheme : lightTheme,
    overrides: pick.value === "dark" ? darkThemeOverrides : lightThemeOverrides,
  }
})


const watermark = ref(false)

let timeout = 1000
const axios = inject("axios");
function loadConsoleServices() {
  axios
    .get(store.base_url + "/console/services")
    .then((res) => {
      store.console_services = res.data;
    })
    .catch((err) => {
      console.error(err);
      setTimeout(() => {
        loadConsoleServices()
      }, timeout < 30000 ? timeout += 1000 : timeout)
    });
}
loadConsoleServices()
</script>

<style>
.n-layout .n-layout-scroll-container {
  overflow-y: hidden
}
</style>