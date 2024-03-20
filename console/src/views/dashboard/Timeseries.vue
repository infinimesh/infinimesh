<template>
    <iframe :src="src" style="width: calc(100% - 5px); height: 90vh; border: none"></iframe>
</template>

<script setup>
import { computed } from "vue"

import { baseURL, useAppStore } from "@/store/app"
import { useNSStore } from "@/store/namespaces"

const as = useAppStore()
const ns = useNSStore()

const src = computed(() => {
    const params = { token: as.token, title: as.me.title, theme: as.theme, api: baseURL, namespace: ns.selected }
    const src = `${as.console_services.timeseries}?a=${btoa(JSON.stringify(params))}`
    return src
})

</script>