<template>
</template>

<script setup>
import { useNSStore } from '@/store/namespaces';
import { useRouter } from "vue-router";

const nss = useNSStore()
const router = useRouter();

let ns = nss.namespaces[nss.selected]
if (!ns || !ns.plugin) {
    console.log("No Namespace selected or Namespace has no Plugins")
    router.push({ name: "Devices" })
} else {
    console.log("Navigating to Main Dashboard")
    router.push({ name: "DashboardMain" })
}
</script>