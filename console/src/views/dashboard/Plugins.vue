<template>
    <n-spin :show="loading" style="min-width: 90%;">
        <n-grid item-responsive y-gap="10">
            <n-grid-item span="24 500:8 600:6 1000:4">
                <n-h1 prefix="bar" align-text type="info">
                    <n-text type="info"> Apps & Plugins </n-text>
                </n-h1>
            </n-grid-item>
            <n-grid-item span="0 600:2 700:4 1000:12 1400:14"> </n-grid-item>
            <n-grid-item span="24 300:12 500:7 600:6 700:5 1000:4 1400:3">
                <n-button strong secondary round type="info" @click="store.fetchPlugins">
                    <template #icon>
                        <n-icon>
                            <refresh-outline />
                        </n-icon>
                    </template>
                    Refresh
                </n-button>
            </n-grid-item>
            <n-grid-item span="24 300:12 500:7 600:6 700:5 1000:4 1400:2" v-if="dev">
                <plugin-create />
            </n-grid-item>
            <n-grid-item span="24">
                <n-alert title="Disclaimer" type="info">
                    <template #icon>
                        <n-icon>
                            <git-network-outline />
                        </n-icon>
                    </template>
                    Some Apps & Plugins are binded to Namespaces, so number of available ones may differ from Namespace
                    to Namespace.
                </n-alert>
            </n-grid-item>
        </n-grid>
        <plugins-pool :plugins="plugins" />
    </n-spin>
</template>

<script setup>
import { watch } from 'vue'
import { NSpin, NGrid, NGridItem, NH1, NText, NButton, NIcon, NSpace, NAlert } from 'naive-ui';
import { RefreshOutline, GitNetworkOutline } from '@vicons/ionicons5';

import { useAppStore } from "@/store/app";
import { usePluginsStore } from "@/store/plugins"
import { useNSStore } from "@/store/namespaces";
import { storeToRefs } from 'pinia';

import PluginsPool from "@/components/plugins/pool.vue"
import PluginCreate from "@/components/plugins/create-drawer.vue"

const as = useAppStore()
const nss = useNSStore();
const store = usePluginsStore()

const { loading, plugins } = storeToRefs(store)
const { dev } = storeToRefs(as)


store.fetchPlugins()

const { selected } = storeToRefs(nss)

watch(selected, () => {
    store.fetchPlugins()
})
</script>