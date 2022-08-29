<template>
  <n-spin :show="loading">
    <n-grid item-responsive y-gap="10" x-gap="10">
      <n-grid-item span="24 500:14 600:12 1000:12">
        <n-space justify="space-between">
          <n-h1 prefix="bar" align-text type="info">
            <n-text type="info"> Devices </n-text>
            (
            <n-number-animation :from="0" :to="devices.length" />
            )
          </n-h1>
          <n-button strong secondary round type="info" @click="handleRefresh">
            <template #icon>
              <n-icon>
                <refresh-outline />
              </n-icon>
            </template>
            Refresh State
          </n-button>
        </n-space>
      </n-grid-item>
      <n-grid-item span="24 300:24 500:10 600:12 1000:12">
        <n-space justify="space-evenly" align="center">
          <device-create />
          <device-register v-if="console_services.handsfree != undefined" />
        </n-space>
      </n-grid-item>
    </n-grid>
    <devices-pool :devices="devices" :show_ns="show_ns" @refresh="() => store.fetchDevices(true, true)" />
  </n-spin>
</template>

<script setup>
import { defineAsyncComponent } from "vue"
import { NSpin, NH1, NText, NIcon, NButton, NGrid, NGridItem, NSpace, NNumberAnimation } from "naive-ui";

import { useAppStore } from "@/store/app";
import { useDevicesStore } from "@/store/devices";
import { storeToRefs } from "pinia";

const RefreshOutline = defineAsyncComponent(() => import("@vicons/ionicons5/RefreshOutline"))
const DevicesPool = defineAsyncComponent(() => import("@/components/devices/pool.vue"))

const DeviceCreate = defineAsyncComponent(() => import("@/components/devices/create-drawer.vue"))
const DeviceRegister = defineAsyncComponent(() => import("@/components/devices/register-modal.vue"))


const store = useDevicesStore();
const { loading, devices_ns_filtered: devices, show_ns } = storeToRefs(store);

store.fetchDevices(true, true);

const { console_services } = storeToRefs(useAppStore())

function handleRefresh() {
  store.fetchDevices(true);
}
</script>