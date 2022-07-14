<template>
    <n-button @click="show = true" type="success" dashed>
        <template #icon>
            <n-icon>
                <keypad-outline />
            </n-icon>
        </template>
        Enter Code
    </n-button>
    <n-modal :show="show" @update:show="e => show = e">
        <n-card style="min-width: 60vw" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <template #header>
                Authorize device(s) with the Code
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
                <n-steps size="small" :current="step" :status="'process'">
                    <n-step title="Pick devices" description="Configure the token your App going to get">
                    </n-step>
                    <n-step title="Enter the Code" description="Enter the code shown in your App">
                    </n-step>
                    <n-step title="Finish" description="Successfuly authorized App">
                    </n-step>
                </n-steps>

                <component :is="content" @update:nextEnabled="v => next_enabled = v" @update:value="value" />

                <n-alert title="Error" type="error" v-if="error">
                    {{ error }}
                </n-alert>
            </n-space>

            <n-space justify="end" align="center" style="margin-top: 2vh">
                <n-button type="error" round secondary @click="show = false">Cancel</n-button>
                <n-button type="success" round @click="next" v-if="step < 3" :disabled="!next_enabled"
                    :loading="loading">Next</n-button>
                <n-button type="info" round @click="show = false" v-else>Close</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>

<script setup>
import { ref, watch, computed, defineAsyncComponent } from "vue"
import { NButton, NIcon, NCard, NModal, NSteps, NStep, NSpace, NAlert } from 'naive-ui';

import { useDevicesStore } from "@/store/devices";
import { useAppStore } from "@/store/app";

const KeypadOutline = defineAsyncComponent(() => import("@vicons/ionicons5/KeypadOutline"))
const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const PickDevice = defineAsyncComponent(() => import("./register_modal/pick.vue"))
const Code = defineAsyncComponent(() => import("./register_modal/code.vue"))

const store = useDevicesStore()
const as = useAppStore()

const show = ref(true)
const step = ref(2)
const next_enabled = ref(false)
const loading = ref(false)
const error = ref(false)

watch(show, () => {
    step.value = 1
    next_enabled.value = false
    error.value = false
})

const content = computed(() => {
    switch (step.value) {
        case 1:
            return PickDevice
        case 2:
            return Code
    }
})

let token = ""
let tokenRequest = {}
let code = ""

function value(val) {
    switch (step.value) {
        case 1:
            tokenRequest = val
            break
        case 2:
            code = val.code
            break
    }
}

async function next() {
    loading.value = true
    try {
        switch (step.value) {
            case 1:
                token = await store.makeDevicesToken(tokenRequest.devices, tokenRequest.post)
                step.value++
                break
            case 2:
                let res = await as.http.post("/handsfree", {
                    code: 1, payload: [code, token]
                })

                if (res.isAxiosError) {
                    throw res
                }
                if (res.data.code == 'SUCCESS') step.value++
                break
        }
        error.value = false
    } catch (e) {
        error.value = e.response.data.message
    }
    loading.value = false
}
</script>