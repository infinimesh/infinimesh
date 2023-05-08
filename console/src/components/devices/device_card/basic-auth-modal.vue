<template>
    <modal-button @submit="handleSubmit" min-width="30vw" submit-text="Close" cancel-text=""
        :type="device.basicEnabled ? 'warning' : 'info'">
        <template #button-text>
            Basic Auth
        </template>
        <template #icon>
            <lock-open-outline v-if="device.basicEnabled" />
            <lock-closed-outline v-else />
        </template>
        <template #header>
            Manage Basic Auth for <strong>{{ device.title }}</strong>
        </template>

        <n-space vertical justify="space-between">

            <n-space justify="start">
                <n-h3 align-text>
                    <n-text type="info">
                        Basic Auth is
                    </n-text>
                </n-h3>

                <n-switch size="large" :rail-style="railStyle" @update:value="emit('toggle')" :value="device.basicEnabled">
                    <template #checked>
                        Enabled
                    </template>
                    <template #checked-icon>
                        <n-icon>
                            <lock-open-outline />
                        </n-icon>
                    </template>
                    <template #unchecked>
                        Disabled
                    </template>
                    <template #unchecked-icon>
                        <n-icon>
                            <lock-closed-outline />
                        </n-icon>
                    </template>
                </n-switch>
            </n-space>

            <n-form label-placement="left">
                <n-form-item label="Username">
                    <n-input-group>
                        <n-input :value="device.uuid" @update:value="" />
                        <n-button type="primary" ghost @click="handleCopy(device.uuid, 'Username')">
                            <template #icon>
                                <n-icon>
                                    <copy-outline />
                                </n-icon>
                            </template>
                        </n-button>
                    </n-input-group>
                </n-form-item>
                <n-form-item label="Password">
                    <n-input-group>
                        <n-input :value="device.certificate.fingerprint" @update:value="" />
                        <n-button type="primary" ghost @click="handleCopy(device.certificate.fingerprint, 'Password')">
                            <template #icon>
                                <n-icon>
                                    <copy-outline />
                                </n-icon>
                            </template>
                        </n-button>
                    </n-input-group>
                </n-form-item>
            </n-form>
        </n-space>
    </modal-button>
</template>

<script setup>
import { ref, toRefs, watch, defineAsyncComponent } from "vue"
import {
    NSpace, NIcon, NSwitch, NH3,
    NText, NForm, NFormItem, NInput,
    NInputGroup, NButton, useMessage
} from "naive-ui";

const ModalButton = defineAsyncComponent(() => import("@/components/core/modal-button.vue"))

const LockClosedOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LockClosedOutline"))
const LockOpenOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LockOpenOutline"))
const CopyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CopyOutline"))

const emit = defineEmits(['toggle'])
const props = defineProps({
    device: {
        required: true
    }
})
const { device } = toRefs(props)

function handleSubmit(close) {
    close()
}

const railStyle = ({
    focused,
    checked
}) => {
    const style = {};
    if (checked) {
        style.background = "#f2c97d";
        if (focused) {
            style.boxShadow = "0 0 0 2px #f2c97d40";
        }
    } else {
        style.background = "#2080f0";
        if (focused) {
            style.boxShadow = "0 0 0 2px #2080f040";
        }
    }
    return style;
}

const msg = useMessage();
async function handleCopy(value, what) {
    try {
        await navigator.clipboard.writeText(value);
        msg.success(`Copied ${what} to clipboard`);
    } catch {
        msg.error(`Failed to copy ${what} to clipboard`);
    }
}
</script>