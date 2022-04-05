<template>
  <n-card hoverable :title="device.title" :header-style="{fontFamily: 'Exo'}" style="border-radius: 0;">
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

    <template #action>
      <n-spin :show="patching">
        <device-state-collapse :state="store.device_state(device.uuid)" :patch="patch" @submit="handlePatchDesired" />
      </n-spin>
      <n-space justify="start" align="center" style="margin-top: 1vh;">
          <n-button
            type="success" round tertiary
            :disabled="subscribed"
            @click="handleSubscribe">{{ subscribed ? 'Subscribed' : 'Subscribe'}}</n-button>
          
          <n-button type="warning" round tertiary @click="patch = !patch">{{ patch ? 'Cancel Patch' : 'Patch Desired' }}</n-button>
      </n-space>
    </template>
  </n-card>
</template>

<script setup>
import { ref, computed } from "vue";
import {
  NCard, NTooltip, NIcon, useMessage, NSpin, useLoadingBar,
  NTag, NSpace, NButton } from "naive-ui"
import { Bulb } from '@vicons/ionicons5'
import DeviceStateCollapse from './state-collapse.vue'

import { useDevicesStore } from "@/store/devices";

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

const store = useDevicesStore()

const subscribed = computed(() => {
  return store.device_subscribed(device.value.uuid)
})

const message = useMessage()
async function handleUUIDClicked() {
  try {
    await navigator.clipboard.writeText(device.value.uuid);
    message.success('Device UUID copied to clipboard')
  } catch {
    message.error('Failed to copy device UUID to clipboard')
  }
}

function handleSubscribe() {
  store.subscribe([device.value.uuid])
}

const bar = useLoadingBar()
const patch = ref(false)
const patching = ref(false)
async function handlePatchDesired(state) {
  console.log(state)
  patching.value = true
  await store.patchDesiredState(device.value.uuid, state, bar)
  patch.value = false
  patching.value = false
}
</script>