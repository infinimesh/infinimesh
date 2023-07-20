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
                <objects-tree :fetch="deletables" @loading="change_loading" @confirm="() => { emit('confirm'); show = false }"/>
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

import {
    NPopconfirm, NModal, NCard,
    NButton, NSpin, NTree,
    NSpace
} from "naive-ui";

import ObjectsTree from "./objects-tree.vue";
import { access_lvl_conv } from "@/utils/access";


const show = ref(false)
const loading = ref(true)

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

function close() {
    loading.value = false
    show.value = false
}
function change_loading(val) {
    console.log('loading', val)
    loading.value = val
}

</script>