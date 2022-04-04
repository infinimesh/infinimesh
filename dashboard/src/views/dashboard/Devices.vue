<template>
  <n-spin :show="loading">
    <devices-pool :devices="devices" />
  </n-spin>
</template>

<script setup>
import { ref } from "vue"
import { NSpin } from "naive-ui"
import { storeToRefs } from "pinia";
import DevicesPool from "@/components/devices/pool.vue"
import { useDevicesStore } from "@/store/devices";

const store = useDevicesStore()
const { loading, devices } = storeToRefs(store)

store.fetchDevices()
</script>