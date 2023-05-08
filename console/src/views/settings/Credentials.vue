<template>
    <n-space>
        <n-h2 prefix="bar">
            <n-text>
                Credentials
            </n-text>
        </n-h2>
    </n-space>

    <n-form ref="form" :model="model" label-placement="left" style="max-width: 640px;">
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

    <n-space justify="start" align="center">
        <n-button type="info" round secondary @click="reset">Reset</n-button>
        <n-button type="warning" round @click="handleSubmit">Submit</n-button>
    </n-space>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent } from "vue";
import { NH2, NText, NForm, NFormItem, NInput, NIcon, NButton, NSpace, useMessage, useLoadingBar, NTabs, NTabPane } from 'naive-ui';

import { useAppStore } from "../../store/app";
import { useAccountsStore } from "@/store/accounts";

const EyeOffOutline = defineAsyncComponent(() => import("@vicons/ionicons5/EyeOffOutline"))
const EyeOutline = defineAsyncComponent(() => import("@vicons/ionicons5/EyeOutline"))

const as = useAppStore();
const store = useAccountsStore();

const model = ref({
    type: 'standard',
    data: [as.me.title, ''],
});

function reset() {
    model.value = {
        type: 'standard',
        data: [as.me.title, ''],
    };

    default_data()
}

function default_data() {
    store.getCredentials(as.me.uuid).then((data) => {
        if (!data || !data.credentials) return
        for (let cred of data.credentials) {
            if (cred.type == model.value.type) {
                model.value.data = cred.data
            }
        }
    })
}

const message = useMessage();
const bar = useLoadingBar();
async function handleSubmit() {
    let err = await store.setCredentials(as.me.uuid, model.value, bar);
    if (!err) {
        message.success('Credentials set');
    } else {
        message.error(`${err.response.status}: ${(err.response.data ?? { message: "Unexpected Error" }).message}`);
    }
}

onMounted(() => {
    default_data()
})

</script>