<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <router-view />
    <n-watermark
      v-if="watermark"
      content="development preview"
      cross
      fullscreen
      :font-size="16"
      :line-height="16"
      :width="250"
      :height="150"
      :x-offset="12"
      :y-offset="80"
      :rotate="-15"
    />
  </n-config-provider>
</template>

<script>
import { computed } from "vue"
import {
  NConfigProvider, NGlobalStyle, NWatermark,
  darkTheme, useOsTheme } from 'naive-ui'

export default {
  components: {
    NConfigProvider, NGlobalStyle, NWatermark
  },
  setup() {
    console.log(darkTheme)
    const osThemeRef = useOsTheme();
    return {
      theme: computed(() => osThemeRef.value === "dark" ? darkTheme : null),
      osTheme: osThemeRef,
      watermark: false,
    };
  }
}
</script>