<template>
  <n-spin :show="loading">
    <n-grid item-responsive y-gap="10" x-gap="10">
      <n-grid-item span="24 500:14 600:12 1000:12">
        <n-space justify="space-between">
          <n-h1 prefix="bar" align-text type="info">
            <n-text type="info"> Devices </n-text>
            (
            <n-number-animation :from="0" :to="filteredDevices.length || 0" />
            )
          </n-h1>
          <n-space>
            <n-dropdown 
              trigger="click"
              width="200px"
              placement="bottom-start"
              :options="filterDeviceOptions"
              @select="handleFilterSelect"
            >
              <n-input
                ref="filterRef"
                placeholder="Filter devices eg. :uuid:abc"
                v-model:value="filterTerm"
                clearable
                @clear="clearFilter"
              />
            </n-dropdown>
          </n-space>
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
    <devices-pool :devices="filteredDevices()" :show_ns="show_ns" @refresh="() => store.fetchDevices(true, true)" />
  </n-spin>
</template>

<script setup>
import { defineAsyncComponent, watch, ref } from "vue"
import { NSpin, NH1, NText, NInput, NIcon, NButton, NGrid, NGridItem, NSpace, NDropdown, NNumberAnimation } from "naive-ui";

import { useAppStore } from "@/store/app";
import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { usePluginsStore } from "@/store/plugins";
import { storeToRefs } from "pinia";

const RefreshOutline = defineAsyncComponent(() => import("@vicons/ionicons5/RefreshOutline"))
const DevicesPool = defineAsyncComponent(() => import("@/components/devices/pool.vue"))

const DeviceCreate = defineAsyncComponent(() => import("@/components/devices/create-drawer.vue"))
const DeviceRegister = defineAsyncComponent(() => import("@/components/devices/register-modal.vue"))


const store = useDevicesStore();
const { loading, devices_ns_filtered: devices, show_ns } = storeToRefs(store);

const filterRef = ref(null);
const filterTerm = ref("");
const filterDeviceOptions = [
  {
    label: ":uuid:",
    key: ":uuid:"
  },
  {
    label: ":enabled:",
    key: ":enabled:"
  },
  {
    label: ":tag:",
    key: ":tag:"
  },
  {
    label: ":title:",
    key: ":title:"
  },
  {
    label: ":namespace:",
    key: ":namespace:"
  }
]

function handleFilterSelect(option) {
  filterTerm.value = filterTerm.value + " " + option;
  filterRef.value.focus();
}

function clearFilter() { filterTerm.value = "" }

function parseFilterText() {
  const filters = {};
  const parts = filterTerm.value.split(' ');
  parts.forEach(part => {
    const [key, value] = part.split(':').filter(String);
    if (key && value) {
      filters[key] = value;
    }
  });

  return filters;
}

function matchFilterDevice(device, filters) {
  for(let key in filters) {
    const filterValue = filters[key];

    if(key === "uuid" && !device.uuid.toLowerCase().includes(filterValue)) {
      return false;
    }
    
    if(key === "enabled") {
      const expectedValue = filterValue === "true";
      if (device.enabled !== expectedValue) {
        return false;
      }
    }

    if(key === "title" && !device.title.toLowerCase().includes(filterValue)) {
      return false;
    }

    if(key === "tag") {
      const filterTags = filterValue.split(",");

      if(!device.tags.some(tag => 
        filterTags.some(filterTag => 
          tag.toLowerCase().includes(filterTag))
        )
      ) {
        return false;
      }
    }

    if(
      key === "namespace" && 
      !device.access.namespace.toLowerCase().includes(filterValue)
    ) {
      return false;
    }
  }

  return true;
}

function filteredDevices() {
  const filters = parseFilterText()
  return devices.value.filter((device) => {
    return matchFilterDevice(device, filters)
  })
}

store.fetchDevices(true, true);

const app = useAppStore()
const { console_services } = storeToRefs(app)

const { selected, namespaces } = storeToRefs(useNSStore())
const plugins = usePluginsStore()

async function load_plugin() {
  plugins.current = false
  if (selected.value == 'all') {
    return
  }

  let ns = namespaces.value[selected.value]

  if (!ns || !ns.plugin) {
    plugins.current = false
    return
  }

  const { data } = await plugins.get(ns.plugin.uuid)
  if (ns.vars) data.vars = ns.vars
  plugins.current = data
}

watch(selected, load_plugin)

function handleRefresh() {
  store.fetchDevices(true);
}

load_plugin()
</script>