<template>
  <n-space justify="space-between" :style="{ padding: '15px' }" align="center">
    <current_thing />

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
  if (route.name == "DashboardMain" && store.dev)
    return h(PluginVars)

  return null
}

function infinimesh() {
  return h('span', { class: "infinimesh" }, [
    "infinimesh", h('span', { class: "copyright" }, "©")
  ])
}

const jollymesh = defineAsyncComponent(() => import("@/assets/icons/jollymesh.svg"))

function current_thing() {
  if (!store.current_thing)
    return infinimesh()

  switch (store.current_thing.k) {
    case 'jolly':
      return h(jollymesh, {
        id: "svgmesh", style: { height: '5vh', cursor: 'pointer' }, onClick: () => store.current_thing.on = !store.current_thing.on
      })
  }

  return infinimesh()

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

#svgmesh text {
  fill: var(--n-text-color);
}
</style>