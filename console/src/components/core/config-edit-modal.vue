<template>
    <modal-button min-width="30vw" submit-text="Submit" cancel-text="Cancel" type="info" @submit="handleSubmit"
        @show="reset">
        <template #icon>
            <cog-outline />
        </template>
        <template #header>
            Manage Configuration Data for <strong>{{ o.title }}</strong>
        </template>

        <n-space vertical justify="space-between">

            <n-input v-model:value="config" type="textarea" placeholder="Desired State" :status="validation" />

        </n-space>
    </modal-button>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, toRefs } from 'vue';

import { NSpace, NInput } from 'naive-ui';

const CogOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CogOutline"))

const ModalButton = defineAsyncComponent(() => import("@/components/core/modal-button.vue"))


const props = defineProps({
    o: {
        type: Object,
        required: true
    }
})

const { o } = toRefs(props)

const emit = defineEmits(["submit"])
const config = ref(o.value.config ? JSON.stringify(o.value.config, null, 2) : '{}');

const validation = computed(() => {
    try {
        let d = JSON.parse(config.value);
        if (typeof d != "object") return "error";
        return "success";
    } catch {
        return "error";
    }
});

function handleSubmit(close) {
    if (validation.value) {
        o.value.config = JSON.parse(config.value);
        emit("submit", o.value)
        close()
    }
}

function reset() {
    config.value = o.value.config ? JSON.stringify(o.value.config, null, 2) : '{}';
}
</script>