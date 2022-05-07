<template>
    <n-tr>
        <n-td colspan="1"></n-td>
        <n-td colspan="1">
            <n-text type="warning" strong>
                List of Accounts, who has access to this Namespace
            </n-text>
        </n-td>
        <n-td colspan="2">
            <n-text type="warning" strong v-if="!loading && joins.length == 0">
                No Accounts found
            </n-text>
        </n-td>
        <n-td>
            <n-button strong quaternary round type="info" @click="load">
                <template #icon>
                    <n-icon>
                        <refresh-outline />
                    </n-icon>
                </template>
                Refresh
            </n-button>
        </n-td>
    </n-tr>
    <n-tr v-if="loading">
        <n-td></n-td>
        <n-td colspan="4" align="start">
            <n-progress type="line" :percentage="100" processing border-radius="12px 0 12px 0" :show-indicator="false"
                fill-border-radius="12px 0 12px 0" style="width: 80%" />
        </n-td>
    </n-tr>
    <n-tr v-for="acc in joins" v-else>
        <n-td></n-td>
        <n-td>
            <uuid-badge :uuid="acc.uuid" />
        </n-td>
        <n-td>
            {{ acc.title }}
        </n-td>
        <n-td>
            <access-badge :access="acc.access.level" join />
            <access-badge access="OWNER" v-if="acc.access.role == 'OWNER'" left="5px" join />
        </n-td>
        <n-td></n-td>
    </n-tr>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { NTr, NTd, NProgress, NText, NButton, NIcon } from "naive-ui"
import { RefreshOutline } from "@vicons/ionicons5";

import UuidBadge from "@/components/core/uuid-badge.vue";
import AccessBadge from "@/components/core/access-badge"

import { useNSStore } from "@/store/namespaces"

const store = useNSStore()

const props = defineProps({
    namespace: {
        type: String,
        required: true
    }
})

const loading = ref(false)
const joins = ref([])

async function load() {
    loading.value = true
    const { data } = await store.loadJoins(props.namespace)
    joins.value = data.accounts
    loading.value = false
}

onMounted(() => {
    load()
})

</script>