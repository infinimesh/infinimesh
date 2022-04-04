<template>
  <n-config-provider :theme="theme" :hljs="hljs">
    <n-message-provider>
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
    </n-message-provider>
  </n-config-provider>
</template>

<script>
import { computed } from "vue"
import {
  NConfigProvider, NGlobalStyle, NWatermark,
  NMessageProvider, darkTheme, useOsTheme } from 'naive-ui'

import hljs from "@/utils/hljs"

export default {
  components: {
    NConfigProvider, NGlobalStyle, NWatermark, NMessageProvider,
  },
  setup() {
    const osThemeRef = useOsTheme();
    return {
      theme: computed(() => osThemeRef.value === "dark" ? darkTheme : null),
      osTheme: osThemeRef,
      watermark: false,
      hljs
    };
  }
}
</script>