<template>
  <n-dropdown :options="options" trigger="hover" @select="v => pick = v">
    <n-button type="primary" size="large" ghost circle><template #icon>
        <n-icon :component="icons[theme]" />
      </template></n-button>
  </n-dropdown>
</template>

<script setup>
import { ref, watch, onMounted } from "vue"
import { NDropdown, NButton, NIcon, useOsTheme } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { useAppStore } from "@/store/app"

import { MoonOutline, SunnyOutline, CogOutline } from '@vicons/ionicons5';

import { renderIcon } from "@/utils";

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