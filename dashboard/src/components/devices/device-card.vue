<template>
  <n-dropdown trigger="contextmenu" :options="options">
    <n-card hoverable :title="device.title" :header-style="{fontFamily: 'Exo'}">
      <template #header-extra>
       <n-tooltip trigger="hover" @click="handleUUIDClicked">
          <template #trigger>
            <span @click="handleUUIDClicked">
              {{ device.uuid_short }}
            </span>
          </template>
          {{ device.uuid }}
        </n-tooltip>
        <n-icon size="2vh" :color="device.enabled ? '#52c41a' : '#eb2f96'" style="margin-left: 1vw;">
          <bulb />
        </n-icon>
      </template>

      <template #footer v-if="device.tags.length > 0">
        Tags:
        <n-tag type="warning" round v-for="tag in device.tags" :key="tag" style="margin-right: 3px">
          {{ tag }}
        </n-tag>
      </template>

      <template #action >
        <n-code
          v-if="state && state.reported && state.reported.version != '0'"
          :code="JSON.stringify(state.reported.data, null, 2)" language="json" />
        <n-code v-else code="No State have been reported yet" />
        <n-collapse>
          <n-collapse-item title="Desired State" name="desired">
            <n-code
              v-if="state && state.desired && state.desired.version != '0'"
              :code="JSON.stringify(state.desired.data, null, 2)" language="json" />
            <n-code v-else code="No Desired state have been set yet" />
          </n-collapse-item>
        </n-collapse>
      </template>
    </n-card>
  </n-dropdown>
</template>

<script setup>
import { ref, computed } from "vue";
import {
  NDropdown, NCard, NTooltip, NIcon, useMessage,
  NTag, NCode, NCollapse, NCollapseItem } from "naive-ui"
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

const message = useMessage()
async function handleUUIDClicked() {
  try {
    await navigator.clipboard.writeText(device.value.uuid);
    message.success('Device UUID copied to clipboard')
  } catch {
    message.error('Failed to copy device UUID to clipboard')
  }
}
</script>