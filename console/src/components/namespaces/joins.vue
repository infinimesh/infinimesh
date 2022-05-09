<template>
    <n-tr>
        <n-td colspan="3" align="center">
            <n-text type="warning" strong>
                List of Accounts, who has access to this Namespace
            </n-text>
        </n-td>
        <n-td colspan="1">
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
            <template v-if="editing == acc.uuid">
                <access-badge access="READ" join :cb="(v) => handleJoin(acc.uuid, v)" />
                <access-badge access="MGMT" join left="5px" :cb="(v) => handleJoin(acc.uuid, v)" />
                <access-badge access="ADMIN" join left="5px" :cb="(v) => handleJoin(acc.uuid, v)" />
                <access-badge access="ROOT" join left="5px" :cb="(v) => handleJoin(acc.uuid, v)" />
            </template>
            <template v-else>
                <access-badge :access="acc.access.level" join />
                <access-badge access="OWNER" v-if="acc.access.role == 'OWNER'" left="5px" join />
            </template>
        </n-td>
        <n-td>
            <n-space>
                <n-button type="success" round secondary @click="() => editing = null" v-if="editing == acc.uuid">
                    Cancel Edit
                </n-button>
                <n-button type="success" round secondary @click="() => editing = acc.uuid" v-else>
                    Change
                </n-button>
                <n-button v-if="admin" type="warning" round secondary @click="handleJoin(acc.uuid, 0)">Remove
                </n-button>
            </n-space>
        </n-td>
    </n-tr>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { NTr, NTd, NProgress, NText, NButton, NIcon, NSpace } from "naive-ui"
import { RefreshOutline } from "@vicons/ionicons5";

import UuidBadge from "@/components/core/uuid-badge.vue";
import AccessBadge from "@/components/core/access-badge"

import { useNSStore } from "@/store/namespaces"
import { access_levels } from "@/utils/access";

const store = useNSStore()

const props = defineProps({
    namespace: {
        type: String,
        required: true
    },
    admin: {
        type: Boolean,
        default: false
    }
})

const loading = ref(false)
const editing = ref(null)
const joins = ref([])

async function load() {
    loading.value = true
    const { data } = await store.loadJoins(props.namespace)
    joins.value = data.accounts
    loading.value = false
}

async function handleJoin(account, access) {
    loading.value = true
    try {
        const { data } = await store.join(props.namespace, account, access_levels[access])
        joins.value = data.accounts
    } catch (e) {
        console.error(e)
    }
    editing.value = null
    loading.value = false
}

onMounted(() => {
    load()
})
</script>