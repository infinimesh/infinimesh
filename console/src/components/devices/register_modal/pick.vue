<template>

  <n-alert :title="post ? 'Posting enabled' : 'Posting disabled'" :type="post ? 'warning' : 'info'">
    {{ post ?
        'App will have enough access to send reports and to desire picked devices state' :
        'App will have access only to read and subscribe to the picked devices state'
    }}
    <n-switch v-model:value="post" />
  </n-alert>

  <n-checkbox-group v-model:value="selected" style="margin-top: 10px">
    <n-table>
      <thead>
        <tr>
          <th>Picked</th>
          <th>UUID</th>
          <th>Title</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="device in devices">
          <td align="center">
            <n-checkbox size="large" class="table-checkbox" style="--n-size: max(32px, min(48px, 2vh)); --n-border-radius: var(--n-size);"
              :value="device.uuid" />
          </td>
          <td>
            <uuid-badge :uuid="device.uuid" :type="device.enabled ? 'success' : 'error'" />
          </td>
          <td>
            {{ device.title }}
          </td>
          <td>
            {{ (namespaces[device.access.namespace] ?? { title: '-' }).title }}
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-checkbox-group>
</template>

<script setup>
import { ref, watch, defineAsyncComponent, computed } from "vue"
import { NTable, NCheckbox, NCheckboxGroup, NSwitch, NAlert } from "naive-ui"

import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { storeToRefs } from "pinia";

const { devices: devicesMap } = storeToRefs(useDevicesStore());
const { namespaces } = storeToRefs(useNSStore());

const devices = computed(() => Object.values(devicesMap.value));

const UuidBadge = defineAsyncComponent(() => import("@/components/core/uuid-badge.vue"))

const post = ref(false)
const selected = ref(null)

const emit = defineEmits(['update:nextEnabled', 'update:value'])

watch(selected, () => {
  if (!selected.value || selected.value.length == 0) emit('update:nextEnabled', false)
  else {
    emit('update:nextEnabled', true)
    emit('update:value', {
      post: post.value,
      devices: selected.value
    })
  }
})

watch(post, () => {
  emit('update:value', {
    post: post.value,
    devices: selected.value
  })
})

</script>
