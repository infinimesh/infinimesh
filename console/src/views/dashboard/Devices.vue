<template>
  <n-spin :show="loading">
    <n-grid item-responsive y-gap="10">
      <n-grid-item span="24 500:8 600:6 1000:4">
        <n-h1 prefix="bar" align-text type="info">
          <n-text type="info"> Devices </n-text>
        </n-h1>
      </n-grid-item>
      <n-grid-item span="0 600:2 700:4 1000:12 1400:14"> </n-grid-item>
      <n-grid-item span="24 300:12 500:7 600:6 700:5 1000:4 1400:3">
        <n-button strong secondary round type="info" @click="emit('refresh')">
          <template #icon>
            <n-icon>
              <refresh-outline />
            </n-icon>
          </template>
          Refresh State
        </n-button>
      </n-grid-item>
      <n-grid-item span="24 300:12 500:7 600:6 700:5 1000:4 1400:2">
        <device-create />
      </n-grid-item>
    </n-grid>
    <devices-pool :devices="devices" :show_ns="show_ns" @refresh="() => store.fetchDevices(true, true)" />
  </n-spin>
</template>

<script setup>
import { NSpin, NH1, NText, NIcon, NButton, NGrid, NGridItem } from "naive-ui";
import { RefreshOutline } from "@vicons/ionicons5";
import { useDevicesStore } from "@/store/devices";
import { storeToRefs } from "pinia";
import DevicesPool from "@/components/devices/pool.vue";
import DeviceCreate from "@/components/devices/create-drawer.vue";

const store = useDevicesStore();
const { loading, devices_ns_filtered: devices, show_ns } = storeToRefs(store);

store.fetchDevices(true, true);
</script>