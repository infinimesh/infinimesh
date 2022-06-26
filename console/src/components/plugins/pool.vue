<template>
    <n-empty size="huge" description="No Plugins found" v-if="plugins.length == 0" style="margin-top: 20vh"></n-empty>
    <div v-else>
        <n-grid cols="1 s:1 m:2 l:3 xl:4 2xl:4" ref="grid" responsive="screen" style="margin-top: 10px">
            <n-grid-item v-for="(col, i) in pool" :key="i">
                <plugin-card v-for="plugin in col" :key="plugin.uuid" :plugin="plugin" />
            </n-grid-item>
        </n-grid>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { NEmpty, NGrid, NGridItem, } from "naive-ui";

import PluginCard from "./plugin-card.vue";

const grid = ref({ responsiveCols: 0 });

const props = defineProps({
    plugins: {
        type: Array,
        required: true,
    },
});

const emit = defineEmits(["refresh"]);

const pool = computed(() => {
    try {
        let plugins = props.plugins;
        let div = (grid.value ?? { responsiveCols: 0 }).responsiveCols;
        if (!div || div == 1) return [plugins];
        let res = new Array(div);
        for (let i = 0; i < div; i++) {
            res[i] = new Array();
        }
        for (let i = 0; i <= plugins.length; i++) {
            for (let j = 0; j < div && i + j < plugins.length; j++) {
                res[j].push(plugins[i + j]);
            }
            i += div - 1;
        }
        return res;
    } catch (e) {
        console.error(e);
        return [];
    }
});
</script>