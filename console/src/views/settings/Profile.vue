<template>
    <n-space>
        <n-h2 prefix="bar">
            <n-text>
                Profile
            </n-text>
        </n-h2>
    </n-space>
    <n-space align="center" @mouseover="hover = true" @mouseleave="hover = false">
        <n-avatar round :size="64" />

        <n-input-group v-if="editing_title">
            <n-input v-model:value="me.title" />
            <n-button type="error" ghost @click="stop_editing">
                <template #icon>
                    <n-icon :component="CloseOutline" />
                </template>
            </n-button>
            <n-button type="success" ghost @click="save_title">
                <template #icon>
                    <n-icon :component="SaveOutline" />
                </template>
            </n-button>
        </n-input-group>
        <n-h1 style="cursor: pointer;" @click="editing_title = true" v-else>
            <n-text>
                {{ me.title }}
            </n-text>
            <n-icon size="20" style="margin-left: 4px;" v-show="hover">
                <pencil-outline />
            </n-icon>
        </n-h1>

        <n-h3 v-show="me.username" style="cursor: pointer;" @click="router.push({ name: 'Credentials' })">
            <n-text :depth="3">
                {{ me.username }}
            </n-text>
        </n-h3>
    </n-space>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent } from 'vue';
import { useRouter } from 'vue-router';

import {
    NSpace, NIcon, NAvatar,
    NH1, NH2, NH3, NText,
    NInputGroup, NInput, NButton,
    useLoadingBar
} from 'naive-ui';

import { useAppStore } from '../../store/app';
import { useAccountsStore } from '../../store/accounts';
import { storeToRefs } from 'pinia';

const PencilOutline = defineAsyncComponent(() => import("@vicons/ionicons5/PencilOutline"))
const SaveOutline = defineAsyncComponent(() => import("@vicons/ionicons5/SaveOutline"))
const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))

const { me } = storeToRefs(useAppStore())
const store = useAccountsStore();
const router = useRouter()

const bar = useLoadingBar()

const hover = ref(false)
const editing_title = ref(false)

async function stop_editing() {
    await store.sync_me()
    editing_title.value = false
}

async function save_title() {
    let res = await store.updateAccount(me.value, bar)
    console.log(res)

    editing_title.value = false
}

onMounted(async () => {
    await store.sync_me()
    store.getCredentials(me.value.uuid).then((data) => {
        if (!data || !data.credentials) return
        for (let cred of data.credentials) {
            if (cred.type == 'standard') {
                me.value.username = cred.data[0]
            }
        }
    })
})
</script>