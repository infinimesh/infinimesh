<template>
    <n-button type="primary" size="large" ghost circle @click="show = true"
        v-show="plugs.current && plugs.current.vars">
        <template #icon>
            <n-icon size="1.8rem">
                <CogOutline />
            </n-icon>
        </template>
    </n-button>
    <n-modal v-model:show="show">
        <n-card style="width: 600px" :bordered="false" size="huge" role="dialog" aria-modal="true"
            :mask-closable="true">
            <template #header>
                Configure <b>{{plugin.title}}</b> plugin variables
            </template>
            <template #header-extra>
                <n-button @click="show = false" quaternary circle size="large">
                    <template #icon>
                        <n-icon>
                            <close-outline />
                        </n-icon>
                    </template>
                </n-button>
            </template>

            <n-table style="margin-top: 2vh; min-width: 50%">
                <thead>
                    <tr>
                        <th>Variable</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="[k, v] in Object.entries(plugin.vars)" :key="k">
                        <td>{{k}}</td>
                        <td>
                            <n-input v-model:value="plugin.vars[k]" type="text" :placeholder="`Enter ${k} value`" />
                        </td>
                    </tr>
                </tbody>
            </n-table>

            <n-space justify="end" align="center" style="margin-top: 2vh">
                <n-button type="error" round secondary @click="show = false">Cancel</n-button>
                <n-button type="warning" round @click="submit">Submit</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>

<script setup>
import { ref } from "vue"

import { NIcon, NButton, NModal, NCard, NTable, NInput, NSpace } from "naive-ui"
import { CogOutline, CloseOutline } from "@vicons/ionicons5"

import { storeToRefs } from "pinia"
import { usePluginsStore } from "@/store/plugins";
import { useNSStore } from "@/store/namespaces"

const show = ref(false)
const plugs = usePluginsStore()
const nss = useNSStore()

const { current: plugin } = storeToRefs(plugs)

async function submit() {
    let ns = nss.namespaces[nss.selected]
    ns.plugin.vars = plugin.value.vars

    await nss.update(ns)
    show.value = false
}
</script>