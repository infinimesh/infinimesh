<template>
    <n-layout>
        <n-layout-header>
            <dashboard-nav :ns-selector="false" />
        </n-layout-header>
        <n-layout has-sider :content-style="{ minHeight: '90vh' }">
            <n-layout-sider collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
                @mouseover="collapsed = false" @mouseleave="collapsed = true">
                <settings-menu :collapsed="collapsed" />
            </n-layout-sider>
            <n-layout-content :content-style="{ padding: noContentPadding ? 0 : '24px' }" :native-scrollbar="false">
                <router-view />
            </n-layout-content>
        </n-layout>
        <n-layout-footer>
            <dashboard-footer />
        </n-layout-footer>
    </n-layout>
</template>

<script setup>
import {
    NLayout,
    NLayoutHeader,
    NLayoutContent,
    NLayoutSider,
    NLayoutFooter,
} from "naive-ui";

import DashboardNav from "@/components/dashboard/nav.vue";
import DashboardFooter from "@/components/core/footer.vue";

import SettingsMenu from "@/components/settings/menu.vue";
</script>