<template>
    <n-button type="info" round tertiary @click="show = true">
        <template #icon>
            <n-icon>
                <add-outline />
            </n-icon>
        </template>
        {{  device.tags.length > 0 ? 'Edit' : 'Add'  }} Tags
    </n-button>
    <n-modal :show="show" @update:show="e => show = e">
        <n-spin :show="loading">
            <template #description>
                Updating Tags...
            </template>
            <n-card style="min-width: 30vw; max-width: 90vw;" :bordered="false" size="huge" role="dialog"
                aria-modal="true">
                <template #header>
                    Edit Device({{  device.title  }}) Tags
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
                    <n-dynamic-tags v-model:value="tags" type="warning" round size="large" />
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
import { useMessage, NButton, NIcon, NCard, NModal, NSpace, NDynamicTags, NSpin } from 'naive-ui';

const AddOutline = defineAsyncComponent(() => import("@vicons/ionicons5/AddOutline"))
const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const props = defineProps({
    device: {
        required: true
    }
})

const { device } = toRefs(props)

const show = ref(false)
const loading = ref(false)

const tags = ref([])
watch(show, () => {
    console.log(show.value, device.value.tags)
    tags.value = show.value ? device.value.tags : []
    loading.value = false
})

const emit = defineEmits(['save'])
const message = useMessage()

function handleSubmit() {
    loading.value = true
    emit('save', tags.value, () => {
        show.value = false
    }, (msg) => {
        message.error(msg)
        loading.value = false
    })
}
</script>