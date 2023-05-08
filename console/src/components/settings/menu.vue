<template>
    <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
        :value="selected" />
</template>

<script setup>
import { ref, h, computed, defineAsyncComponent } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { NMenu } from "naive-ui";

import { renderIcon } from "@/utils";

import { useAppStore } from "@/store/app"
import { storeToRefs } from "pinia"

const Person = defineAsyncComponent(() => import("@vicons/ionicons5/Person"));
const KeyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/KeyOutline"));
const LockClosedOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LockClosedOutline"));


const props = defineProps({
    collapsed: {
        type: Boolean,
        default: false,
    },
});

const route = useRoute();
const selected = computed(() => route.name);

function renderLabelLink(route, label = false) {
    if (!label) {
        label = route
    }
    return () => h(
        RouterLink,
        {
            to: {
                name: route,
            },
        },
        { default: () => label }
    )
}

const menuOptions = ref([
    {
        label: renderLabelLink("Profile"),
        key: "Profile",
        icon: renderIcon(Person),
    },
    {
        label: renderLabelLink("Credentials"),
        key: "Credentials",
        icon: renderIcon(LockClosedOutline),
    },
    {
        label: renderLabelLink("Tokens"),
        key: "Tokens",
        icon: renderIcon(KeyOutline),
    }
])
</script>