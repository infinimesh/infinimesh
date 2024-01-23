<template>
    <iframe v-if="plugin && src != ''" style="width: calc(100% - 5px); height: 90vh; border: none" :src="src" allow="clipboard-write"></iframe>
    <n-space align="center" justify="center" class="fullscreen padded" style="min-height: 60vh;">

        <n-alert title="No Namespace Selected" type="info" v-if="!ns">
            <template #icon>
                <n-icon>
                    <git-network-outline />
                </n-icon>
            </template>
            Dashboard Plugins are binded to the namespace, so you must select a namespace to view Dashboards.
            <n-space justify="center" style="margin-top: 20px">
                <ns-selector />
            </n-space>
        </n-alert>

        <n-alert title="Namespace has no Plugin connected" type="info" v-else-if="!ns.plugin">
            <template #icon>
                <n-icon>
                    <extension-puzzle-outline />
                </n-icon>
            </template>
            Navigate to Apps & Plugins page in the menu to pick a Plugin for your project.
            <n-space justify="center" style="margin-top: 20px">
                <n-button type="success" ghost @click="router.push({ name: 'Plugins' })">
                    <template #icon>
                        <n-icon>
                            <open-outline />
                        </n-icon>
                    </template>
                    Pick a Plugin
                </n-button>
            </n-space>
        </n-alert>

        <n-alert title="Coudln't get plugin :(" type="error" style="min-width: 40vh;" v-else-if="state == 'notfound'">
            <template #icon>
                <n-icon>
                    <extension-puzzle-outline />
                </n-icon>
            </template>
            <ul>
                <li>
                    Make sure plugin still exists
                </li>
                <li>
                    Make sure plugin provider is online
                </li>
                <li>
                    Contact tech support
                </li>
            </ul>
        </n-alert>

        <n-alert title="In a few..." type="info" style="min-width: 30vw" v-else-if="state == 'loading'">
            <template #icon>
                <n-spin>
                    <n-icon>
                        <log-in-outline />
                    </n-icon>
                </n-spin>
            </template>
            Loading plugin
        </n-alert>

        <n-alert :title="`Namespace has ${plugin.kind} Plugin connected`" type="warning"
            v-else-if="state == 'wrongkind'">
            <template #icon>
                <n-icon>
                    <extension-puzzle-outline />
                </n-icon>
            </template>
            Plugin assosiated with this Namespace is not compatible with Dashboard page
            <n-space justify="center" style="margin-top: 20px">
                <ns-selector />
            </n-space>
        </n-alert>
    </n-space>
</template>

<script setup>
import { ref, computed, watch, onMounted, defineAsyncComponent } from "vue"
import { useRouter } from "vue-router";

import { NSpace, NAlert, NIcon, NButton, NSpin } from "naive-ui";
import { GitNetworkOutline, ExtensionPuzzleOutline, OpenOutline } from "@vicons/ionicons5";

import NsSelector from "@/components/core/ns-selector.vue";
import { useNSStore } from '@/store/namespaces';
import { usePluginsStore } from "@/store/plugins";
import { baseURL, useAppStore } from "@/store/app";

const as = useAppStore()
const nss = useNSStore()
const plugs = usePluginsStore()
const router = useRouter();

const ns = computed(() => nss.namespaces[nss.selected])
const plugin = ref(false)
const state = ref(false)

async function loadPlugin() {
    state.value = "loading"
    plugin.value = false

    if (!ns.value || !ns.value.plugin || !ns.value.plugin.uuid) {
        plugs.current = false
        state.value = 'notfound'
        return
    }

    let { uuid, vars } = ns.value.plugin

    try {
        const data = await plugs.get(uuid)
        if (vars) data.vars = vars
        plugs.current = data
        plugin.value = data

        if (data.kind != "EMBEDDED") {
            state.value = "wrongkind"
            return
        }
        state.value = false

    } catch (e) {
        state.value = 'notfound'
        plugs.current = false
    }
}

onMounted(async () => {
    await loadPlugin()
    watch(ns, loadPlugin)
})

const src = computed(() => {
    console.log(state.value, !plugin.value, plugin.value.kind != "EMBEDDED", state.value || !plugin.value || plugin.kind != "EMBEDDED")
    if (state.value || !plugin.value || plugin.value.kind != "EMBEDDED") return ""

    const params = { token: as.token, title: as.me.title, namespace: nss.selected, theme: as.theme, api: baseURL, vars: plugin.value.vars }
    const src = `${plugin.value.embeddedConf.frameUrl}?a=${btoa(JSON.stringify(params))}`
    return src
})

const LogInOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LogInOutline"))
</script>

<style scoped>
.padded {
    padding-top: 24px
}
</style>
