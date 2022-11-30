<template>
  <n-space justify="space-between" align="center">
    <n-spin size="small" v-show="loading" />
    <n-select v-model:value="selected" :options="options" :style="{ minWidth: '15vw' }" />
    <n-button type="primary" size="small" ghost circle @click="() => store.fetchNamespaces()">
      <template #icon>
        <n-icon>
          <refresh-outline />
        </n-icon>
      </template>
    </n-button>
  </n-space>
</template>

<script setup>
import { computed, watch, defineAsyncComponent } from "vue";
import { NSpace, NSpin, NSelect, NIcon, NButton } from "naive-ui";
import { useNSStore } from "@/store/namespaces";
import { storeToRefs } from "pinia";

const RefreshOutline = defineAsyncComponent(() => import("@vicons/ionicons5/RefreshOutline"))

const store = useNSStore();

function shortUUID(uuid) {
  return uuid.substr(0, 8);
}

const { loading, selected, namespaces_list: namespaces } = storeToRefs(store);
const options = computed(() => {
  return [
    { label: "All", value: "all" },
    ...namespaces.value.map((ns) => ({
      label: `${ns.title} (${shortUUID(ns.uuid)})`,
      value: ns.uuid,
    })),
  ];
});

watch(namespaces, () => {
  if (!selected.value) {
    selected.value = (namespaces.value[0] ?? { uuid: "all" }).uuid;
  }
});

store.fetchNamespaces();
</script>