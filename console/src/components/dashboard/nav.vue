<template>
  <n-space justify="space-between" :style="{ padding: '15px' }" align="center">
    <span class="infinimesh">infinimesh<span class="copyright">©</span></span>
    <ns-selector />
    <div style="margin-left: 5px">
      <n-space justify="space-between">
        <theme-picker />
        <router-action />
      </n-space>
    </div>
    <user-details />
  </n-space>
</template>

<script setup>
import { h, defineAsyncComponent } from "vue"
import { useRoute } from "vue-router"
import { NSpace } from "naive-ui";

import { useAppStore } from "@/store/app"

const NsSelector = defineAsyncComponent(() => import("@/components/core/ns-selector.vue"))
const ThemePicker = defineAsyncComponent(() => import("@/components/core/theme-picker.vue"))
const UserDetails = defineAsyncComponent(() => import("@/components/dashboard/user.vue"))

const PluginVars = defineAsyncComponent(() => import("@/components/plugins/plugin-vars.vue"))

const store = useAppStore()
const route = useRoute()

function RouterAction() {
  console.log(route.name)
  console.log(store.dev)
  if (route.name == "DashboardMain" && store.dev)
    return h(PluginVars)

  return null
}
</script>

<style>
span.infinimesh {
  font-size: 3vh;
  font-family: "Exo 2", sans-serif;
}

span.copyright {
  vertical-align: super;
  font-size: 1.2vh;
}
</style>