<template>
  <n-empty size="huge" description="No Devices" v-if="devices.length == 0"></n-empty>
  <n-grid v-else cols="1 s:1 m:2 l:3 xl:4 2xl:4" responsive="screen" ref="grid">
    <n-grid-item v-for="(col, i) in pool" :key="i">
      <device-card v-for="device in col" :key="device.uuid" :device="device" :show_ns="show_ns" />
    </n-grid-item>
  </n-grid>
</template>

<script setup>
import { ref, computed } from "vue";
import { NEmpty, NGrid, NGridItem } from "naive-ui";
import DeviceCard from "./device-card.vue"

const grid = ref({ responsiveCols: 0 })

const props = defineProps({
  devices: {
    type: Array,
    required: true,
  },
  show_ns: {
    type: Boolean,
    default: false,
  }
})

const pool = computed(() => {
  try {
    let devices = props.devices
    let div = (grid.value ?? { responsiveCols: 0 }).responsiveCols
    if (!div || div == 1) return [devices]
    let res = new Array(div);
    for (let i = 0; i < div; i++) {
      res[i] = new Array();
    }
    for (let i = 0; i <= devices.length; i++) {
      for (let j = 0; j < div && i + j < devices.length; j++) {
        res[j].push(devices[i + j]);
      }
      i += div - 1;
    }
    return res;
  } catch (e) {
    console.error(e)
    return []
  }
})
</script>