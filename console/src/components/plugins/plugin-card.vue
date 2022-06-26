<template>
    <n-spin :show="patching">

        <n-card hoverable :title="plugin.title" :header-style="{ fontFamily: 'Exo 2' }" style="border-radius: 0">
            <template #header-extra>
                <n-tooltip trigger="hover" @click="handleUUIDClicked">
                    <template #trigger>
                        <n-tag
                            :color="{ textColor: kinds[plugin.kind].color ?? '#52c41a', borderColor: kinds[plugin.kind].color ?? '#52c41a' }"
                            size="large" round>
                            {{ kinds[plugin.kind].label }}
                        </n-tag>
                    </template>
                    {{ kinds[plugin.kind].desc }}
                </n-tooltip>
                <n-tooltip trigger="hover" @click="handleUUIDClicked">
                    <template #trigger>
                        <n-tag :color="{ textColor: '#52c41a', borderColor: '#52c41a' }" style="margin-left: 5px"
                            size="large" round @click="handleUUIDClicked">
                            {{ plugin.uuid_short }}
                        </n-tag>
                    </template>
                    {{ plugin.uuid }}
                </n-tooltip>
            </template>

            <template #cover>
                <div style="max-width: 90%; max-height: 30vh" v-if="plugin.logo">
                    <img :src="plugin.logo" style="padding: 20px;">
                </div>
                <n-space align="center" justify="center" v-else>
                    <n-icon size="15vh">
                        <image-outline />
                    </n-icon>
                </n-space>
            </template>

            <template #footer>
                <vue-markdown-it :source='plugin.description' />
            </template>

            <template #action v-if="!props.preview">
                <n-space justify="start">
                    <n-button strong secondary round type="success" v-if="plugin.kind != 'UNKNOWN'">
                        <template #icon>
                            <n-icon>
                                <add-outline />
                            </n-icon>
                        </template>
                        Use it with your Namespace
                    </n-button>

                    <n-popconfirm @positive-click="handleDelete">
                        <template #trigger>
                            <n-button v-if="dev" type="error" round secondary>
                                Delete
                            </n-button>
                        </template>
                        Are you sure about deleting this plugin?
                    </n-popconfirm>
                </n-space>
            </template>
        </n-card>
    </n-spin>
</template>

<script setup>
import { ref, computed } from "vue";
import {
    NCard,
    NTooltip,
    NIcon,
    useMessage,
    NSpin,
    NTag,
    NSpace,
    NButton,
    NPopconfirm,
} from "naive-ui";
import { AddOutline, ImageOutline } from "@vicons/ionicons5"

import VueMarkdownIt from 'vue3-markdown-it';

import { useAppStore } from "@/store/app";
import { usePluginsStore } from "@/store/plugins";
import { storeToRefs } from "pinia";

const as = useAppStore()
const { dev } = storeToRefs(as)

const props = defineProps({
    plugin: {
        type: Object,
        required: true,
    },
    show_ns: {
        type: Boolean,
        default: false,
    },
    preview: {
        type: Boolean,
        default: false
    }
});

const plugin = computed(() => {
    let r = props.plugin;
    if (props.preview) {
        r.uuid = "infinimesh will give some unique ID to this plugin and it will be shown here"
        r.uuid_short = "s0m3uu1d"
        return r
    }
    r.uuid_short = r.uuid.substr(0, 8);
    return r;
});

const message = useMessage();
async function handleUUIDClicked() {
    try {
        await navigator.clipboard.writeText(device.value.uuid);
        message.success("Device UUID copied to clipboard");
    } catch {
        message.error("Failed to copy device UUID to clipboard");
    }
}

const kinds = {
    UNKNOWN: {
        label: "UNKNOWN",
        desc: "This App or Plugin will not work, please contact to Platform administrators",
        color: "#fc1703"
    },
    EMBEDDED: {
        label: "Embedded",
        desc: "This Plugin will embed into this Console as the main Dashboard",
        color: "#8a2be2"
    }
}

const store = usePluginsStore()
const patching = ref(false)
async function handleDelete() {
    patching.value = true
    try {
        await store.delete(plugin.value.uuid)
        store.fetchPlugins()
    } catch (e) {
        console.error(e)
    }
    patching.value = false
}

</script>