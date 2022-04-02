<template>
  <n-space justify="space-between" align="center">
    <n-spin size="small" v-show="loading" />
    <n-select v-model:value="selected" :options="options" :style="{minWidth: '15vw'}" @update:show="handleShow" />
  </n-space>
</template>

<script setup>
import { ref, computed, watch } from "vue"
import { NSpace, NSpin, NSelect } from "naive-ui"
import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces"
import { storeToRefs } from "pinia";

const store = useNSStore()

const loading = ref(false)

const { selected, namespaces } = storeToRefs(store)
const options = computed(() => {
  return namespaces.value.map(ns => ({
    label: ns.title,
    value: ns.uuid,
  }))
})

const axios = useAppStore().http
async function loadNamespaces() {
  loading.value = true
  axios.get('http://localhost:8000/namespaces', {}).then(res => {
    namespaces.value = res.data.namespaces
    loading.value = false
  })
}

function handleShow(show) {
  if (show) {
    loadNamespaces()
  }
}

watch(namespaces, () => {
  if (!selected.value) {
    selected.value = namespaces.value[0].uuid
  }
})

loadNamespaces()
</script>