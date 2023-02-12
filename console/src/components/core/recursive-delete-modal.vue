<template>
    <n-popconfirm @positive-click="() => show = true">
        <template #trigger>
            <n-button v-if="o.access.role == 'OWNER' || access_lvl_conv(o) > 3" type="error" round secondary
                @click.stop.prevent>Delete
            </n-button>
        </template>
        <span>
            Are you sure about deleting {{ type }} <b>{{ o.title }}</b>?
        </span>
    </n-popconfirm>
    <n-modal v-model:show="show">
        <n-spin :show="loading">
            <n-card style="width: 600px" title="To be deleted" :bordered="false" size="huge" role="dialog"
                aria-modal="true">
                <n-tree :data="tree" :default-expand-all="true" />
                <template #footer>
                    <n-space justify="end">
                        <n-button type="info" round secondary @click="close">
                            Cancel
                        </n-button>
                        <n-button type="error" round secondary @click="() => { emit('confirm'); show = false }">
                            Confirm Delete
                        </n-button>
                    </n-space>
                </template>
            </n-card>
        </n-spin>
    </n-modal>
</template>

<script setup>
import { h, ref, watch } from "vue"

import { NPopconfirm, NModal, NCard, NButton, NSpin, NTree, NSpace, useMessage } from "naive-ui";
import { useAccountsStore } from "@/store/accounts";
import { useNSStore } from "@/store/namespaces";
import { useDevicesStore } from "@/store/devices";

import { access_lvl_conv } from "@/utils/access";

const accs = useAccountsStore();
const nss = useNSStore()
const devs = useDevicesStore();

const show = ref(false)
const loading = ref(true)
const tree = ref([])

const { o, type, deletables } = defineProps({
    o: {
        type: Object,
        required: true
    },
    deletables: {
        default: () => async () => [],
    },
    type: {
        type: String,
        default: "namespace"
    }
})

const emit = defineEmits(['confirm'])

const message = useMessage()
watch(show, async (val) => {
    if (val) {
        loading.value = true

        try {
            const { data } = await deletables()

            if (data.nodes.length == 1) {
                close()
                emit('confirm')
                return
            }

            if (!Object.keys(accs.accounts).length) {
                accs.fetchAccounts()
            }
            if (!Object.keys(devs.devices).length) {
                devs.fetchDevices(false)
            }

            tree.value = makeTree(data.nodes)
            loading.value = false

        } catch (e) {
            message.error("Error loading deletables: " + e.response.statusText)
            close()
        }
    }
})

function close() {
    loading.value = false
    show.value = false
}

function makeTree(data, parent = '') {
    let nodes = []
    data = data.filter((value) => {
        if (value.node == "") {
            return false
        }
        if (value.parent == parent) {
            nodes.push(value)
            return false
        }
        return true
    })

    nodes = nodes.map(node => {
        return {
            key: node.node,
            label: () => resolve(node.node),
            children: makeTree(data, node.node)
        }
    })
    return nodes.length ? nodes : null
}

function resolve(node) {
    const [id, uuid] = node.split('/')
    let type, title;
    switch (id) {
        case "Accounts":
            type = "Account"
            title = accs.accounts[uuid]?.title
            break
        case "Namespaces":
            type = "Namespace"
            title = nss.namespaces[uuid]?.title
            break
        case "Devices":
            type = "Device"
            title = devs.devices[uuid]?.title
            break
        default:
            type = id.slice(0, id.length)
            title = ""
    }

    return h('span', [type + " ", h('b', title), `(${uuid})`])
}

</script>