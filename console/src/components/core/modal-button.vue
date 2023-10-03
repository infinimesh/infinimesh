<template>
    <n-button :type="type" :round="!!$slots['button-text']" :circle="!$slots['button-text']" tertiary @click="show = true">
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
            <n-card :style="{ minWidth, maxWidth, width }" :bordered="false" size="huge" role="dialog"
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
                    <n-button type="error" round secondary @click="handleCancel" v-if="cancelText != ''">{{ cancelText }}</n-button>
                    <n-button type="success" round @click="handleSubmit" :disabled="submitDisabled">{{
                        submitText
                    }}</n-button>
                </n-space>
            </n-card>
        </n-spin>
    </n-modal>
</template>

<script setup>
import { ref, defineAsyncComponent, watch } from "vue";
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
    },
    width: {
        type: String,
        required: false,
        default: null
    },
    minWidth: {
        type: String,
        required: false,
        default: "30vw"
    },
    maxWidth: {
        type: String,
        required: false,
        default: "90vw"
    }
})

const emit = defineEmits(['submit', 'cancel', 'show'])

const show = ref(false)
watch(show, (v) => {
    emit('show', v)
})

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