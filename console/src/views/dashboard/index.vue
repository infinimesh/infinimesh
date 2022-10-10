<template>
    <n-space v-if="!ns" align="center" justify="center" class="fullscreen padded">
        <n-alert title="No Namespace Selected" type="info">
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
    </n-space>
    <n-space v-else-if="!ns.plugin" align="center" justify="center" class="fullscreen padded">
        <n-alert title="Namespace has no Plugin connected" type="info">
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
    </n-space>
    <n-spin :show="plugin.state == 'loading'" style="min-width: 90%;" v-else>
        <n-space v-if="plugin.state == 'notfound'" align="center" justify="center" class="fullscreen padded">
            <n-alert title="Coudln't get plugin :(" type="error" style="min-width: 40vh;">
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
        </n-space>
        <iframe v-else style="width: calc(100% - 5px); height: 90vh; border: none" :src="src"></iframe>
    </n-spin>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue"
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
const plugin = ref({ state: "loading" })

async function loadPlugin() {
    if (!ns.value.plugin || !ns.value.plugin.uuid) return

    let uuid = ns.value.plugin.uuid
    try {
        const { data } = await plugs.get(uuid)
        plugin.value = data
    } catch (e) {
        plugin.value = { state: 'notfound' }
    }
}

onMounted(async () => {
    await loadPlugin()
    watch(ns, loadPlugin)
})

function makeid(length) {
    var result = '';
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() *
            charactersLength));
    }
    return result;
}

const src = computed(() => {
    if (plugin.value.state != undefined) return ""
    const params = { token: as.token, title: as.me.title, namespace: nss.selected, theme: as.theme, api: baseURL }
    const src = `${plugin.value.embeddedConf.frameUrl}?a=${btoa(JSON.stringify(params))}&${makeid(6)}=${makeid(10)}`
    return src
})
</script>

<style scoped>
.padded {
    padding-top: 24px
}
</style>