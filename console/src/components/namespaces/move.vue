<template>
    <modal-button type="warning" @submit="handleSubmit">
        <template #button-text>
            Move
        </template>
        <template #loading-text>
            Moving {{ type }} to the new namespace
        </template>
        <template #header>
            Move {{ obj.title }}({{ obj.uuid.substr(0, 8) }})
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
    </modal-button>
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
const ModalButton = defineAsyncComponent(() => import("@/components/core/modal-button.vue"))

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

const loading = ref(false)
const err = ref()

const emit = defineEmits(['move'])
function handleSubmit(close) {
    err.value = false
    loading.value = true

    emit('move', ns.value, () => {
        err.value = false
        loading.value = false
        close()
    }, (msg) => {
        err.value = msg
        loading.value = false
    })
}

</script>