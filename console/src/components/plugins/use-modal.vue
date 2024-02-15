<template>
    <n-modal :show="props.show">
        <n-card style="width: 600px" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <template #header>
                Use <b>{{ props.plugin.title }}</b> with your Project
            </template>
            <template #header-extra>
                <n-button @click="emit('close')" quaternary circle size="large">
                    <template #icon>
                        <n-icon>
                            <close-outline />
                        </n-icon>
                    </template>
                </n-button>
            </template>

            <n-select v-model:value="namespace" :options="namespaces" :style="{ minWidth: '15vw' }" filterable />

            <n-alert :title="validator_state.title" :type="validator_state.type" v-if="validator_state"
                style="margin-top: 2vh">
                <template #icon>
                    <n-icon>
                        <git-network-outline v-if="validator_state.icon == 'namespace'" />
                        <extension-puzzle-outline v-if="validator_state.icon == 'plugin'" />
                    </n-icon>
                </template>
                {{ validator_state.message }}
            </n-alert>

            <n-table v-if="vars" style="margin-top: 2vh; min-width: 50%">
                <thead>
                    <tr>
                        <th>Variable</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="k in Object.keys(vars)" :key="k">
                        <td>{{ k }}</td>
                        <td>
                            <n-input v-model:value="vars[k]" type="text" :placeholder="`Enter ${k} value`" />
                        </td>
                    </tr>
                </tbody>
            </n-table>

            <n-space justify="end" align="center" style="margin-top: 2vh">
                <n-button type="error" round secondary @click="emit('close')">Cancel</n-button>
                <n-button type="warning" round @click="submit">Submit</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>

<script setup>
import { ref, computed, watch, defineAsyncComponent, toRefs } from 'vue';

import { NCard, NModal, NButton, NIcon, NSelect, NAlert, NSpace, NTable, NInput } from "naive-ui"

import { useNSStore } from "@/store/namespaces";
import { access_lvl_conv } from "@/utils/access";

const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))
const GitNetworkOutline = defineAsyncComponent(() => import("@vicons/ionicons5/GitNetworkOutline"))
const ExtensionPuzzleOutline = defineAsyncComponent(() => import("@vicons/ionicons5/ExtensionPuzzleOutline"))

const nss = useNSStore();

function shortUUID(uuid) {
    return uuid.substr(0, 8);
}

const namespaces = computed(() => {
    return nss.namespaces_list.filter(ns => access_lvl_conv(ns) > 2).map((ns) => ({
        label: `${ns.title} (${shortUUID(ns.uuid)})`,
        value: ns.uuid,
    }));
});

const emit = defineEmits(['close'])
const props = defineProps({
    show: {
        type: Boolean
    },
    plugin: {
        type: Object
    }
})

const { show } = toRefs(props)

const namespace = ref(nss.selected)
const vars = ref(false)

watch(show, (show) => {
    if (!show) {
        vars.value = false
        return
    }

    const plugin = props.plugin
    if (!plugin || !plugin.vars || plugin.vars.length == 0) { vars.value = false; return };

    vars.value = plugin.vars.reduce((r, v) => {
        r[v] = ''; return r
    }, {})
})


const validator_state = computed(() => {
    if (namespace.value == "all") {
        return {
            icon: "namespace",
            title: 'Snap!',
            type: "error",
            message: "You must select particular Namespace"
        }
    }
    let ns = nss.namespaces[namespace.value]
    if (ns.plugin) {
        return {
            icon: "plugin",
            title: "Warning!",
            type: "warning",
            message: "Namespace already has plugin, submitting this will overwrite it!"
        }
    }
    return null
})

async function submit() {
    let ns = nss.namespaces[namespace.value]
    if (!ns) {
        return
    }

    let vars_raw = vars.value
    if (!vars_raw) {
        vars_raw = {}
    }

    await nss.update({ ...ns, plugin: { uuid: props.plugin.uuid, vars: vars_raw } })

    emit('close')
}
</script>