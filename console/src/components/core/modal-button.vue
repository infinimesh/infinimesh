<template>
    <n-button :type="type" round tertiary @click="show = true">
        <template #icon v-if="$slots.icon">
            <n-icon>
                <slot name="icon"></slot>
            </n-icon>
        </template>
        <slot name="button-text"></slot>
    </n-button>
    <n-modal :show="show" @update:show="e => show = e">
        <n-spin :show="loading">
            <template #description>
                <slot name="loading-text" v-if="$slots['loading-text']"></slot>
                <template v-else>
                    Loading...
                </template>
            </template>
            <n-card style="min-width: 30vw; max-width: 90vw;" :bordered="false" size="huge" role="dialog"
                aria-modal="true">

                <template #header>
                    <slot name="header"></slot>
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


                <slot name="default"></slot>

                <n-space justify="end" align="center" style="margin-top: 2vh">
                    <n-button type="error" round secondary @click="handleCancel">{{ cancelText }}</n-button>
                    <n-button type="success" round @click="handleSubmit" :disabled="submitDisabled">{{
                        submitText
                    }}</n-button>
                </n-space>
            </n-card>
        </n-spin>
    </n-modal>
</template>

<script setup>
import { ref, defineAsyncComponent } from "vue";
import {
    NButton, NModal, NCard, NSpin,
    NIcon, NSpace
} from "naive-ui"

const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const props = defineProps({
    type: {
        type: String,
        required: false,
        default: "info"
    },
    cancelText: {
        type: String,
        required: false,
        default: "Cancel"
    },
    submitText: {
        type: String,
        required: false,
        default: "Submit"
    },
    submitDisabled: {
        type: Boolean,
        required: false,
        default: false
    },
    loading: {
        type: Boolean,
        required: false,
        default: false
    }
})

const show = ref(false)

const emit = defineEmits(['submit', 'cancel'])

function handleCancel() {
    show.value = false
    emit('cancel')
}
function handleSubmit() {
    emit('submit', () => {
        show.value = false
    })
}
</script>