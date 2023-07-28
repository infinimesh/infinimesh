<template>
    <iframe :src="src" style="width: calc(100% - 5px); height: 90vh; border: none"></iframe>
</template>

<script setup>
import { computed } from "vue"

import { baseURL, useAppStore } from "@/store/app"

const as = useAppStore()

const src = computed(() => {
    const params = { token: as.token, title: as.me.title, theme: as.theme, api: baseURL }
    const src = `${as.console_services.chatting}?a=${btoa(JSON.stringify(params))}`
    return src
})

</script>