<template>

  <n-config-provider :theme="theme" :theme-overrides="overrides" :hljs="hljs">
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
import { ref, computed } from "vue";
import {
  NConfigProvider,
  NGlobalStyle,
  NWatermark,
  NLoadingBarProvider,
  NMessageProvider,
  darkTheme,
  lightTheme,
  useOsTheme,
} from "naive-ui";

import lightThemeOverrides from "@/assets/light-theme-overrides.json"
import darkThemeOverrides from "@/assets/dark-theme-overrides.json"

import hljs from "@/utils/hljs";

const osThemeRef = useOsTheme();
const theme = computed(() => (osThemeRef.value === "dark" ? darkTheme : lightTheme))
const overrides = computed(() => (osThemeRef.value === "dark" ? darkThemeOverrides : lightThemeOverrides))

const watermark = ref(false)
</script>

<style>
.n-layout .n-layout-scroll-container {
  overflow-y: hidden
}
</style>