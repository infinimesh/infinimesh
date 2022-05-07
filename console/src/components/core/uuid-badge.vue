<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button tertiary round :type="type" @click.stop.prevent="handleCopy()">
                <template #icon>
                    <n-icon>
                        <copy-outline />
                    </n-icon>
                </template>
                {{ shortUUID(uuid) }}
            </n-button>
        </template>
        {{ uuid }}
    </n-tooltip>
</template>

<script setup>
import { NTooltip, NIcon, NButton, useMessage } from 'naive-ui';
import { CopyOutline } from '@vicons/ionicons5';

const { uuid, type } = defineProps({
    uuid: {
        type: String,
        required: true
    },
    type: {
        type: String,
        default: 'info'
    }
})

function shortUUID(uuid) {
    return uuid.substr(0, 8);
}

const message = useMessage();
async function handleCopy() {
    try {
        await navigator.clipboard.writeText(uuid);
        message.success("UUID Copied to clipboard");
    } catch {
        message.error("Failed to copy UUID to clipboard");
    }
}
</script>