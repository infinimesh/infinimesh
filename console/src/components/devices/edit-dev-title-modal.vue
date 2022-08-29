<template>
    <n-button type="info" round quaternary @click="show = true">
        <template #icon>
            <n-icon>
                <pencil-outline />
            </n-icon>
        </template>
    </n-button>
    <n-modal :show="show" @update:show="e => show = e">
        <n-spin :show="loading">
            <template #description>
                Updating Title...
            </template>
            <n-card style="min-width: 30vw; max-width: 90vw;" :bordered="false" size="huge" role="dialog"
                aria-modal="true">
                <template #header>
                    Edit Device({{  device.title  }}) Title
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
                    <n-input v-model:value="title" placeholder="Make it bright" />
                </n-space>

                <n-space justify="end" align="center" style="margin-top: 2vh">
                    <n-button type="error" round secondary @click="show = false">Cancel</n-button>
                    <n-button type="success" round @click="handleSubmit">Submit</n-button>
                </n-space>
            </n-card>
        </n-spin>
    </n-modal>
</template>

<script setup>
import { ref, toRefs, watch, defineAsyncComponent } from "vue"
import { useMessage, NButton, NIcon, NCard, NModal, NSpace, NInput, NSpin } from 'naive-ui';

const PencilOutline = defineAsyncComponent(() => import("@vicons/ionicons5/PencilOutline"))
const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const props = defineProps({
    device: {
        required: true
    }
})

const { device } = toRefs(props)

const show = ref(false)
const loading = ref(false)
const title = ref("")

watch(show, () => {
    title.value = show.value ? device.value.title : ""
    loading.value = false
})

const emit = defineEmits(['save'])
const message = useMessage()

function handleSubmit() {
    loading.value = true
    emit('save', title.value, () => {
        show.value = false
    }, (msg) => {
        message.error(msg)
        loading.value = false
    })
}
</script>