<template>
  <n-spin :show="loading">
    <devices-pool
      :devices="devices"
      :show_ns="show_ns"
      @refresh="() => store.fetchDevices()"
    />
  </n-spin>
</template>

<script setup>
import { NSpin } from "naive-ui";
import { useDevicesStore } from "@/store/devices";
import { storeToRefs } from "pinia";
import DevicesPool from "@/components/devices/pool.vue";

const store = useDevicesStore();
const { loading, devices_ns_filtered: devices, show_ns } = storeToRefs(store);

store.fetchDevices();
</script>