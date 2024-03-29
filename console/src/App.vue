<template>
  <n-config-provider :theme="theme.it" :theme-overrides="theme.overrides" :hljs="hljs">
    <n-loading-bar-provider>
      <n-notification-provider>
        <n-message-provider>
          <n-global-style />
          <router-view />
          <n-watermark v-if="dev" content="dev mode" cross fullscreen :font-size="16" :line-height="16" :width="250"
            :height="150" :x-offset="12" :y-offset="80" :rotate="-15" />

          <current_thing />

        </n-message-provider>
      </n-notification-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup>
import { computed, inject, defineAsyncComponent, watchEffect, h, onMounted } from "vue";
import {
  NConfigProvider,
  NGlobalStyle,
  NWatermark,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  darkTheme,
  lightTheme,
} from "naive-ui";

import { storeToRefs } from 'pinia';
import { useAppStore } from "@/store/app"

import hljs from "@/utils/hljs";

const store = useAppStore()
const { theme: pick, dev } = storeToRefs(store)
const theme = computed(() => {
  return {
    it: pick.value === "dark" ? darkTheme : lightTheme,
    overrides: pick.value === "dark" ? DarkThemeOverrides : LightThemeOverrides,
  }
})


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

onMounted(() => {
  let now = new Date()
  let m = now.getMonth() + 1, d = now.getDate()

  let todo = false

  if ((m == 12) || (m == 1 && d < 14)) { // from 20.12 till 14.01
    todo = {
      k: 'jolly', on: true
    }
  }

  if (todo && store.current_thing && !store.current_thing.on) return

  store.current_thing = todo
})

const snowflakes = defineAsyncComponent(() => import("@/components/core/snowflakes.vue"))

function current_thing() {
  if (!store.current_thing) h('')

  switch (store.current_thing.k) {
    case 'jolly':
      if (store.current_thing.on) return h(snowflakes)
      break
  }

  return null
}

</script>

<style>
.n-layout .n-layout-scroll-container {
  overflow-y: hidden
}
</style>