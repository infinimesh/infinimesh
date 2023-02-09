<template>
    <n-button type="warning" round tertiary @click="show = true">
        Move
    </n-button>
    <n-modal :show="show" @update:show="e => show = e">
        <n-spin :show="loading">
            <template #description>
                Moving {{ type }} to the new namespace
            </template>
            <n-card style="min-width: 30vw; max-width: 90vw;" :bordered="false" size="huge" role="dialog"
                aria-modal="true">
                <template #header>
                    Move {{ obj.title }}({{ obj.uuid.substr(0, 8) }})
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

                <n-space vertical justify="space-between">
                    <span>
                        Current Namespace: <strong>{{
                            nss.namespaces[obj.access.namespace]?.title || obj.access.namespace
                        }}</strong>
                    </span>

                    <n-space align="center">
                        <span>
                            Move to:
                        </span>

                        <n-select v-model:value="ns" :options="namespaces"
                            :style="{ minWidth: '15vw', display: 'inline-block' }" filterable />
                    </n-space>

                    <n-alert type="error" title="Snap!" v-if="err">
                        {{ err }}
                    </n-alert>
                </n-space>

                <n-space justify="end" align="center" style="margin-top: 2vh">
                    <n-button type="error" round secondary @click="show = false">Cancel</n-button>
                    <n-button type="success" round @click="handleSubmit" :disabled="!ns">Submit</n-button>
                </n-space>
            </n-card>
        </n-spin>
    </n-modal>
</template>

<script setup>
import { ref, computed, defineAsyncComponent } from "vue"
import {
    NButton, NModal, NCard,
    NSpace, NSpin, NIcon,
    NSelect, NAlert
} from "naive-ui"
import { useNSStore } from "@/store/namespaces";
import { access_lvl_conv } from "@/utils/access";

const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const props = defineProps({
    type: {
        required: true
    },
    obj: {
        required: true
    }
})

const ns = ref()
const nss = useNSStore()
const namespaces = computed(() => {
    return nss.namespaces_list.map((ns) => ({
        label: `${ns.title} (${ns.uuid.substr(0, 8)})`,
        value: ns.uuid,
        disabled: access_lvl_conv(ns) < 3 || ns.uuid == props.obj.access?.namespace
    })).sort((a, b) => a.disabled - b.disabled);
});

const show = ref(false)
const loading = ref(false)
const err = ref()

const emit = defineEmits(['move'])
function handleSubmit() {
    err.value = false
    loading.value = true

    emit('move', ns.value, () => {
        show.value = false
        err.value = false
        loading.value = false
    }, (msg) => {
        err.value = msg
        loading.value = false
    })
}

</script>