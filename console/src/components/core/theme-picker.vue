<template>
  <n-dropdown :options="options" trigger="hover" @select="v => pick = v">
    <n-button type="primary" size="large" ghost circle><template #icon>
        <n-icon :component="icons[theme]" />
      </template></n-button>
  </n-dropdown>
</template>

<script setup>
import { ref, watch, onMounted, defineAsyncComponent } from "vue"
import { NDropdown, NButton, NIcon, useOsTheme } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { useAppStore } from "@/store/app"

import { renderIcon } from "@/utils";

const MoonOutline = defineAsyncComponent(() => import("@vicons/ionicons5/MoonOutline"))
const SunnyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/SunnyOutline"))
const CogOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CogOutline"))

const { theme, theme_pick: pick } = storeToRefs(useAppStore())

const icons = {
  'light': SunnyOutline,
  'dark': MoonOutline,
}

const options = [{
  key: 'system',
  label: 'System',
  icon: renderIcon(CogOutline),
}, {
  key: 'dark',
  label: 'Dark',
  icon: renderIcon(MoonOutline),
}, {
  key: 'light',
  label: 'Light',
  icon: renderIcon(SunnyOutline),
}]

const osThemeRef = useOsTheme();

function setTheme(val) {
  theme.value = val === "system" ? osThemeRef.value : pick.value
}

watch(pick, setTheme)
onMounted(() => setTheme(pick.value))
</script>