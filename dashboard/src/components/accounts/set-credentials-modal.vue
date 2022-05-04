<template>
    <n-modal v-model:show="props.show" @mask-click="emit('close')" @esc="emit('close')">
        <n-card style="width: 600px" title="Manage Credentials" :bordered="false" size="huge" role="dialog"
            aria-modal="true">
            <template #header-extra>
                <n-button @click="emit('close')" quaternary circle size="large">
                    <template #icon>

                        <n-icon>
                            <close-outline />
                        </n-icon>
                    </template>
                </n-button>
            </template>
            <n-form ref="form" :model="model" label-placement="top">
                Credentials:
                <n-tabs v-model="model.type">
                    <n-tab-pane name="standard" display-directive="if" tab="Standard(user/pass)">
                        <n-form-item label="Username" path="data[0]">
                            <n-input v-model:value="model.data[0]" />
                        </n-form-item>
                        <n-form-item label="Password" path="data[1]">
                            <n-input v-model:value="model.data[1]" type="password" show-password-on="click">
                                <template #password-visible-icon>
                                    <n-icon :size="16" :component="EyeOffOutline" />
                                </template>
                                <template #password-invisible-icon>
                                    <n-icon :size="16" :component="EyeOutline" />
                                </template>
                            </n-input>
                        </n-form-item>
                    </n-tab-pane>
                </n-tabs>
            </n-form>

            <template #footer>
                <n-space justify="end" align="center">
                    <n-button type="error" round secondary @click="emit('close')">Cancel</n-button>
                    <n-button type="info" round secondary @click="reset">Reset</n-button>
                    <n-button type="warning" round @click="handleSubmit">Submit</n-button>
                </n-space>
            </template>
        </n-card>
    </n-modal>
</template>

<script setup>
import { ref, watch } from "vue";
import { NModal, NCard, NForm, NFormItem, NInput, NIcon, NButton, NSpace, useMessage, useLoadingBar, NTabs, NTabPane } from 'naive-ui';
import { CloseOutline, EyeOffOutline, EyeOutline } from '@vicons/ionicons5';

import { useAccountsStore } from "@/store/accounts";

const store = useAccountsStore();

const props = defineProps({
    show: {
        type: Boolean,
        default: false,
    },
    account: {
        type: Object,
        default: () => ({}),
    },
});
const emit = defineEmits(['close'])

const model = ref({
    type: 'standard',
    data: [props.account.title, ''],
});

function reset() {
    model.value = {
        type: 'standard',
        data: [props.account.title, ''],
    };
}

watch(
    () => props.show,
    (val) => {
        reset();
    }
);

const message = useMessage();
const bar = useLoadingBar();
async function handleSubmit() {
    let err = await store.setCredentials(props.account.uuid, model.value, bar);
    if (!err) {
        message.success('Credentials set');
        emit('close');
    } else {
        message.error(`${err.response.status}: ${(err.response.data ?? { message: "Unexpected Error" }).message}`);
    }
}

</script>