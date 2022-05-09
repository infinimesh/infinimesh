<template>
    <n-input-group v-if="show">
        <n-input :status="status" placeholder="Namespace title" v-model:value="title" />
        <n-button-group>
            <n-button type="success" dashed @click="handleCreate">
                <template #icon>
                    <n-icon>
                        <checkmark-outline />
                    </n-icon>
                </template>
            </n-button>
            <n-button type="error" dashed @click="show = false">
                <template #icon>
                    <n-icon>
                        <ban-outline />
                    </n-icon>
                </template>
            </n-button>
        </n-button-group>
    </n-input-group>
    <n-button @click="show = true" type="success" dashed v-else>
        <template #icon>
            <n-icon>
                <add-outline />
            </n-icon>
        </template>
        Create Namespace
    </n-button>
</template>

<script setup>
import { ref } from "vue"

import { NButtonGroup, NButton, NIcon, NInputGroup, NInput } from 'naive-ui';
import { AddOutline, CheckmarkOutline, BanOutline } from '@vicons/ionicons5';
import { useNSStore } from "@/store/namespaces";

const show = ref(false)
const title = ref("")
const status = ref(undefined)

const store = useNSStore()
async function handleCreate() {
    if (!title.value || title.value == "") {
        status.value = "error"
        return
    }
    status.value = "success"

    store.loading = true
    try {
        await store.create({ title: title.value })
        status.value = undefined
        show.value = false
        store.fetchNamespaces()
    } catch {
        status.value = "error"
        store.loading = false
    }
}
</script>