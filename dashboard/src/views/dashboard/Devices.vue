<template>
  <h1>I am Devices Dashboard</h1>
</template>

<script setup>
import { ref } from "vue"
import { useAppStore } from "@/store/app";

const loading = ref(false)

const axios = useAppStore().http
async function loadDevices() {
  loading.value = true
  axios.get('http://localhost:8000/devices', {}).then(res => {
    console.log(res.data)
    loading.value = false
  })
}

loadDevices()
</script>