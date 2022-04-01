<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <router-view />
  </n-config-provider>
</template>

<script>
import { computed } from "vue"
import {
  NConfigProvider, NGlobalStyle,
  darkTheme, useOsTheme } from 'naive-ui'

export default {
  components: {
    NConfigProvider, NGlobalStyle
  },
  setup() {
    console.log(darkTheme)
    const osThemeRef = useOsTheme();
    return {
      theme: computed(() => osThemeRef.value === "dark" ? darkTheme : null),
      osTheme: osThemeRef
    };
  }
}
</script>