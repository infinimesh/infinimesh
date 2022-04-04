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

      <template #footer v-if="device.tags.length > 0">
        Tags:
        <n-tag type="warning" round v-for="tag in device.tags" :key="tag" style="margin-right: 3px">
          {{ tag }}
        </n-tag>
      </template>

      <template #action v-if="state && state.reported && state.reported.version != '0'">
        <n-code :code="JSON.stringify(state.reported.data, null, 2)" language="json" />
      </template>
    </n-card>
  </n-dropdown>
</template>

<script setup>
import { ref, computed } from "vue";
import { NDropdown, NCard, NTooltip, NIcon, NTag, NCode } from "naive-ui"
import { OpenOutline, Bulb } from '@vicons/ionicons5'

import { useDevicesStore } from "@/store/devices";

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

const store = useDevicesStore()
const state = computed(() => store.device_state(device.value.uuid))
</script>