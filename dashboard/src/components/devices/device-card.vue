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

      <template #action>
        <n-collapse v-model:expanded-names="expanded">
          <n-collapse-item title="Reported State" name="reported">
            <n-code :code="reported ? JSON.stringify(reported.data, null, 2) : '// No State have been reported yet'" language="json" />
          </n-collapse-item>
          <n-space justify="space-between" align="center" v-if="reported && expanded.includes('reported')">
            <n-statistic label="Version">
              <n-number-animation
                :from="0"
                :to="reported.version"
                :active="true"
              />
            </n-statistic>
            <n-statistic label="Timestamp">
              <n-date-picker input-readonly :value="(new Date(reported.timestamp)).getTime()" type="datetime" disabled class="pseudo-disabled" />
            </n-statistic>
          </n-space>
          <n-collapse-item title="Desired State" name="desired">
            <n-code :code="desired ? JSON.stringify(desired.data, null, 2) : '// No Desired state have been set yet'" language="json" />
          </n-collapse-item>
          <n-space justify="space-between" align="center" v-if="desired && expanded.includes('desired')">
            <n-statistic label="Version">
              <n-number-animation
                :from="0"
                :to="desired.version"
                :active="true"
              />
            </n-statistic>
            <n-statistic label="Timestamp">
              <n-date-picker input-readonly :value="(new Date(desired.timestamp)).getTime()" type="datetime" disabled class="pseudo-disabled" />
            </n-statistic>
          </n-space>
        </n-collapse>
        <n-space justify="start" align="center" style="margin-top: 1vh;">
            <n-button
              type="success" round tertiary
              :disabled="subscribed"
              @click="handleSubscribe">{{ subscribed ? 'Subscribed' : 'Subscribe'}}</n-button>
            <n-button type="warning" round tertiary @click="handlePatchDesired">Patch</n-button>
        </n-space>
      </template>
    </n-card>
  </n-dropdown>
</template>

<script setup>
import { ref, computed } from "vue";
import {
  NDropdown, NCard, NTooltip, NIcon, useMessage,
  NTag, NCode, NCollapse, NCollapseItem, NDatePicker,
  NSpace, NButton, NStatistic, NNumberAnimation } from "naive-ui"
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

const expanded = ref(['reported'])

const store = useDevicesStore()

const reported = computed(() => {
  let state = store.device_state(device.value.uuid)
  if (!state || !state.reported || state.reported.version == '0') {
    return false
  }
  return state.reported
})

const desired = computed(() => {
  let state = store.device_state(device.value.uuid)
  if (!state || !state.desired || state.desired.version == '0') {
    return false
  }
  return state.desired
})

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

function handlePatchDesired() {
  if(!expanded.value.includes('desired')) {
    expanded.value.push('desired')
  }
}
</script>