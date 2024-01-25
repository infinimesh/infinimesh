<template>
    <modal-button :loading="loading" @submit="handleSubmit" @show="handleShow" min-width="50vw" submit-text="Close" cancel-text="">
        <template #button-text>
            Share
        </template>
        <template #icon>
            <share-social-outline />
        </template>
        <template #loading-text>
            {{ loading_text }}
        </template>
        <template #header>
            Share Device <strong>{{ device.title }}</strong>
        </template>

        <n-space vertical justify="space-between">
            <n-space justify="space-between">
                <n-h3 prefix="bar" align-text>
                    <n-text type="info">
                        Shared to:
                    </n-text>
                </n-h3>

                <n-button strong secondary round type="info" @click="fetch">
                    <template #icon>
                        <n-icon>
                            <refresh-outline />
                        </n-icon>
                    </template>
                    Refresh
                </n-button>

                <n-button strong secondary round type="success" @click="add = !add ? {} : false">
                    <template #icon>
                        <n-icon>
                            <close-outline v-if="add" />
                            <share-social-outline v-else />
                        </n-icon>
                    </template>
                    {{ add? 'Cancel': 'Share' }}
                </n-button>
            </n-space>

            <n-table>
                <thead>
                    <tr>
                        <th>
                            Who
                        </th>
                        <th>
                            Access
                        </th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="add">
                        <td style="width: 40%"><n-select v-model:value="add.uuid" :options="accounts" filterable /></td>
                        <td style="width: 40%">
                            <access-badge :disabled="add.access == 'READ'" access="READ" join
                                :cb="(v) => add.access = v" device />
                            <access-badge :disabled="add.access == 'MGMT'" access="MGMT" join left="5px"
                                :cb="(v) => add.access = v" device />
                        </td>
                        <td>
                            <n-button strong quaternary round type="warning" @click="handleAdd">
                                <template #icon>
                                    <n-icon>
                                        <add-outline />
                                    </n-icon>
                                </template>
                                Submit
                            </n-button>
                        </td>
                    </tr>
                    <render-row v-for="node in pool.slice(10 * page, 10 * page + 10)" :node="node" />
                </tbody>
            </n-table>
            <n-empty size="huge" description="No Entries found" v-if="pool.length == 0"></n-empty>
            <n-pagination v-model:page="page" :page-count="pool.length / 10" />
        </n-space>
    </modal-button>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, h } from "vue"
import {
    NSpace, NH3, NText, NEmpty,
    NButton, NIcon, NTable, NTooltip,
    NPagination, NSelect, useMessage
} from "naive-ui";
import { useDevicesStore } from "@/store/devices";
import { useAccountsStore } from "@/store/accounts"
import { Level } from "infinimesh-proto/build/es/node/access/access_pb"

const ModalButton = defineAsyncComponent(() => import("@/components/core/modal-button.vue"))
const AccessBadge = defineAsyncComponent(() => import("@/components/core/access-badge"))

const AddOutline = defineAsyncComponent(() => import("@vicons/ionicons5/AddOutline"))
const ShareSocialOutline = defineAsyncComponent(() => import("@vicons/ionicons5/ShareSocialOutline"))
const RefreshOutline = defineAsyncComponent(() => import("@vicons/ionicons5/RefreshOutline"))
const CloseOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloseOutline"))
const TrashOutline = defineAsyncComponent(() => import("@vicons/ionicons5/TrashOutline"))

const ObjectIcons = {
    Accounts: defineAsyncComponent(() => import("@vicons/ionicons5/PeopleOutline")),
    Namespaces: defineAsyncComponent(() => import("@vicons/ionicons5/GitNetworkOutline"))
}
const UnknownType = defineAsyncComponent(() => import("@vicons/ionicons5/ExtensionPuzzleOutline"))

const loading = ref(false)
const loading_text = ref("")

const add = ref(false)

const props = defineProps({
    device: {
        type: Object, required: true
    }
})

const pool = ref([])
const page = ref(0)

const store = useDevicesStore()

async function fetch() {
    loading_text.value = "Fetching Records..."
    loading.value = true

    try {
        let { nodes } = await store.fetchJoins(props.device.uuid)
        console.log(nodes)
        pool.value = nodes
    } catch { }

    loading.value = false
    loading_text.value = ""
}

function handleShow(show) {
    if (show) fetch()
}

function handleSubmit(close) {
    close()
}

const accs = useAccountsStore()
if (Object.keys(accs.accounts).length == 0) {
    accs.fetchAccounts()
}

const accounts = computed(() => Object.values(accs.accounts).map(acc => {
    return {
        label: `${acc.title} (${acc.uuid.substr(0, 8)})`,
        value: acc.uuid,
    }
}))


const message = useMessage()
async function handleAdd() {
    loading_text.value = "Adding Record..."
    loading.value = true

    try {
        await store.join({
            node: props.device.uuid,
            join: `Accounts/${add.value.uuid}`,
            access: Level[add.value.access]
        })
        add.value = false
    } catch (e) {
        loading.value = false
        loading_text.value = ""

        message.error("Couldn't add Share Record: " + e.message)
    }
    fetch()
}
async function handleDelete(uuid) {
    loading_text.value = "Removing Record..."
    loading.value = true

    try {
        await store.join({
            node: props.device.uuid,
            join: `Accounts/${uuid}`,
            access: 0
        })
    } catch (e) {
        loading.value = false
        loading_text.value = ""

        message.error("Couldn't remove Share Record: " + e.message)
    }
    fetch()
}

function RenderRow({ node }) {
    let [type, uuid] = node.node.split('/')
    console.log(type, uuid)
    return h('tr', [
        h('td', h(NSpace, [
            h(NIcon, { size: 20 }, () => h(ObjectIcons[type] ?? UnknownType)),
            h(
                NTooltip,
                {
                    trigger: "hover",
                    placement: "top",
                },
                {
                    trigger: () => h('strong', accs.accounts[uuid].title),
                    default: () => uuid
                }
            )
        ])),
        h('td', h(AccessBadge, { access: node.access.level, device: true })),
        h('td', h(NButton, {
            type: 'warning', strong: true,
            secondary: true, round: true,
            onClick: () => handleDelete(uuid)
        }, {
            icon: h(NIcon, () => h(TrashOutline)),
            default: "Delete"
        }))
    ])
}
</script>

<style>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>