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

const { loading, selected, namespaces } = storeToRefs(store)
const options = computed(() => {
  return [
    { label: "All", value: "all" },
    ...namespaces.value.map(ns => ({
      label: ns.title,
      value: ns.uuid,
    }))
  ]
})

function handleShow(show) {
  if (show) {
    store.fetchNamespaces()
  }
}

watch(namespaces, () => {
  if (!selected.value) {
    selected.value = namespaces.value[0].uuid
  }
})

store.fetchNamespaces()
</script>