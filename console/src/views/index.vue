<template>
</template>

<script setup>
import { useRouter } from "vue-router";

import { useAppStore } from '@/store/app';
import { useNSStore } from '@/store/namespaces';

import { useNotification } from "naive-ui"

const store = useAppStore()
const nss = useNSStore()
const router = useRouter();

let ns = nss.namespaces[nss.selected]
if (!ns || !ns.plugin || ns.plugin.kind != "EMBEDDED") {
    console.log("No Namespace selected or Namespace has no Plugins")
    router.push({ name: "Devices" })
} else {
    console.log("Navigating to Main Dashboard")
    router.push({ name: "DashboardMain" })
}

store.notify = useNotification()
</script>