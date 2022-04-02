<template>
  <n-menu :options="menuOptions" :value="selected" />
</template>

<script setup>
import { ref, h, computed } from "vue"
import { RouterLink, useRoute } from "vue-router"
import { NMenu } from "naive-ui"

import { renderIcon } from "@/utils"
import { HardwareChipOutline } from "@vicons/ionicons5"

const route = useRoute()
const selected = computed(() => route.name)
const menuOptions = ref([
  {
    label: () => h(RouterLink, {
      to: {
        name: "Devices",
      }
    }, { default: () => "Devices" }),
    key: "Devices",
    icon: renderIcon(HardwareChipOutline)
  }
])

</script>