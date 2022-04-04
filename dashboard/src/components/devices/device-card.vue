<template>
  <n-dropdown trigger="contextmenu" :options="options">
    <n-card hoverable :title="device.title" :header-style="{fontFamily: 'Exo'}">
      <template #header-extra>
       <n-tooltip trigger="hover">
          <template #trigger>
            {{ device.uuid_short }}
          </template>
          {{ device.uuid }}
        </n-tooltip>
        <n-icon size="2vh" :color="device.enabled ? '#52c41a' : '#eb2f96'">
          <bulb />
        </n-icon>
      </template>
    </n-card>
  </n-dropdown>
</template>

<script setup>
import { ref, computed } from "vue";
import { NDropdown, NCard, NTooltip, NIcon } from "naive-ui"
import { OpenOutline, Bulb } from '@vicons/ionicons5'

import { renderIcon } from "@/utils";

const props = defineProps({
  device: {
    type: Object,
    required: true,
  },
})

const device = computed(() => {
  let r = props.device
  r.uuid_short = r.uuid.substr(0, 8)
  return r
})
const options = ref([
  {
    key: 'open', label: 'Open',
    icon: renderIcon(OpenOutline),
    props: {
      onClick: () => {
        console.log('open')
      }
    }
  }
])
</script>