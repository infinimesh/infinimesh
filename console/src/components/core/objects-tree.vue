<template>
    <n-tree :data="tree" :default-expand-all="true" />
</template>

<script setup>
import { ref, h, watch, onMounted, defineAsyncComponent } from "vue";

import {
    NTree, NIcon, useMessage
} from "naive-ui"

import { useAccountsStore } from "@/store/accounts";
import { useNSStore } from "@/store/namespaces";
import { useDevicesStore } from "@/store/devices";

const ObjectIcons = {
    Devices: defineAsyncComponent(() => import("@vicons/ionicons5/HardwareChipOutline")),
    Accounts: defineAsyncComponent(() => import("@vicons/ionicons5/PeopleOutline")),
    Namespaces: defineAsyncComponent(() => import("@vicons/ionicons5/GitNetworkOutline"))
}
const UnknownType = defineAsyncComponent(() => import("@vicons/ionicons5/ExtensionPuzzleOutline"))

const accs = useAccountsStore();
const nss = useNSStore()
const devs = useDevicesStore();

const tree = ref([])

const { fetch } = defineProps({
    fetch: {
        default: () => async () => [],
    }
})

const emit = defineEmits(['loading', 'confirm'])

const message = useMessage()
async function load() {
    emit('loading', true);
    try {
        const { data } = await fetch()

        if (data.nodes.length == 1) {
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
    } catch (e) {
        console.log(e)
        message.error("Error loading fetch: " + e.response.statusText)
    }
    emit('loading', false);
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
        const [id, uuid] = node.node.split('/')

        return {
            key: node.node,
            label: () => resolve(id, uuid),
            prefix: () => h(NIcon, null, () => h(ObjectIcons[id] ?? UnknownType)),
            children: makeTree(data, node.node)
        }
    })
    return nodes.length ? nodes : null
}

function resolve(id, uuid) {
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

onMounted(() => {
    load()
})
</script>