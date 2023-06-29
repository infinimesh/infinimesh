<template>
    <n-space vertical justify="start" align="start">
        <n-h2 prefix="bar">
            <n-text>
                Tokens
            </n-text>
        </n-h2>

        <n-space style="padding: 10px; border: 1px solid #333; border-radius: 10px;" vertical>
            <n-h3>
                <n-text>
                    Create Personal Access Token
                </n-text>
            </n-h3>

            <n-space align="center">
                <n-text>
                    Expire At:
                </n-text>
                <n-date-picker v-model:value="expire_at" type="datetime"
                    :is-date-disabled="ts => ts <= new Date().getTime()" :is-time-disabled="timeDisabled" />
                <n-dropdown trigger="hover" :options="options" @select="handleSelectPreset">
                    <n-button>Presets</n-button>
                </n-dropdown>
            </n-space>

            <n-form label-placement="left">
                <n-form-item>
                    <template #label>
                        Client
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :component="HelpCircleOutline" />
                            </template>

                            Who or What will be using this token? Recommended template is {app} | {device} | {os}
                        </n-tooltip>
                    </template>
                    <n-input v-model:value="client" />
                </n-form-item>
            </n-form>

            <n-divider dashed />
            <n-text>
                Confirm your password to continue.
            </n-text>
            <n-form inline v-if="!pat">
                <n-form-item label="Username">
                    <n-input v-model:value="credentials[0]" />
                </n-form-item>
                <n-form-item label="Password">
                    <n-input v-model:value="credentials[1]" type="password" />
                </n-form-item>
            </n-form>

            <n-button type="info" ghost :loading="pat_loading" @click="handleGeneratePAT" :disabled="!credentials[1]">
                Generate
            </n-button>

            <n-input-group style="margin-top: 10px; max-width: 512px;" v-if="pat">
                <n-input :value="pat" @update:value="" />
                <n-button type="info" @click="handleCopy">
                    <template #icon>
                        <n-icon :component="CopyOutline" />
                    </template>
                </n-button>
            </n-input-group>
        </n-space>

        <sessions />
    </n-space>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent } from "vue"

import {
    NSpace, NH2, NH3, NText,
    NDatePicker, NDropdown,
    NButton, NInput, NInputGroup,
    NIcon, useMessage, NForm,
    NFormItem, NDivider, NTooltip
} from 'naive-ui';

import { useAppStore } from "../../store/app"
import { useAccountsStore } from "../../store/accounts";

import { UAParser } from "ua-parser-js"

import sessions from "@/components/settings/sessions.vue";

const CopyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CopyOutline"))
const HelpCircleOutline = defineAsyncComponent(() => import("@vicons/ionicons5/HelpCircleOutline"))

const expire_at = ref(new Date().getTime() + 86400000)
const pat = ref("")
const pat_loading = ref(false)

const as = useAppStore()
const store = useAccountsStore()

const credentials = ref([as.me.username, ""])

let res = {};
  try {
    res = new UAParser(navigator.userAgent).getResult()
  } catch (e) {
    console.warn("Failed to get user agent", e)
  }

const client = ref(`Console | ${res.os?.name ?? 'Unknown'} | ${res.browser?.name ?? 'Unknown'}`)

function timeDisabled(ts) {
    return {
        isHourDisabled: (hour) => {
            return ts < new Date().getTime
        }
    }
}

const options = [
    {
        label: '1 Day',
        key: 86400000,
    },
    {
        label: '1 Week',
        key: 604800000,
    },
    {
        label: '1 Month',
        key: 2592000000,
    },
    {
        label: '3 Month',
        key: 7776000000,
    },
    {
        label: '1 Year',
        key: 31536000000,
    },
]

function handleSelectPreset(key) {
    expire_at.value = new Date().getTime() + key
}

const msg = useMessage()
async function handleGeneratePAT() {
    pat_loading.value = true

    as.http.post(as.base_url + "/token", {
        auth: {
            type: 'standard',
            data: credentials.value
        },
        exp: Math.round(expire_at.value / 1000),
        client: client.value
    })
        .then((res) => {
            pat.value = res.data.token;
            msg.success('Token generated successfully!')
        })
        .catch((e) => {
            msg.error(e.response.data.message)
        }).finally(() => {
            pat_loading.value = false
        });
}

onMounted(async () => {
    await store.sync_me()
    store.getCredentials(as.me.uuid).then((data) => {
        if (!data || !data.credentials) return
        for (let cred of data.credentials) {
            if (cred.type == 'standard') {
                as.me.username = cred.data[0]
            }
        }
    })
})

async function handleCopy() {
    try {
        await navigator.clipboard.writeText(pat.value);
        msg.success("Token Copied to clipboard");
    } catch {
        msg.error("Failed to copy Token to clipboard");
    }
}

</script>