<template>
    <n-layout>
        <n-layout-header>
            <n-space justify="space-between" :style="{ padding: '15px' }" align="center">
                <n-skeleton height="40px" width="10vw" round />

                <n-skeleton height="40px" width="15vw" :sharp="false" />

                <n-skeleton height="40px" width="40px" round />

                <n-skeleton height="40px" width="8vw" :sharp="false" />
            </n-space>
        </n-layout-header>
        <n-layout has-sider :content-style="{ minHeight: '90vh' }">
            <n-layout-sider :width="240" :collapsed="false">
                <n-space justify="space-around" align="center" style="margin-top: 30px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
                <n-space justify="space-around" align="center" style="margin-top: 6px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
                <n-space justify="space-around" align="center" style="margin-top: 6px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
                <n-space justify="space-around" align="center" style="margin-top: 6px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
                <n-space justify="space-around" align="center" style="margin-top: 6px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
                <n-space justify="space-around" align="center" style="margin-top: 6px">
                    <n-skeleton height="30px" width="30px" round />
                    <n-skeleton text :repeat="1" style="width: 140px" />
                </n-space>
            </n-layout-sider>
            <n-layout-content :content-style="{ padding: noContentPadding ? 0 : '24px' }" :native-scrollbar="false">
                <n-space justify="center">
                    <n-space vertical align="center">
                        <n-icon color="#70c0e8ff" size="250">
                            <cloud-offline-outline />
                        </n-icon>
                        <n-h1>
                            <n-text type="primary">
                                Offline :(
                            </n-text>
                        </n-h1>
                        <n-h4>
                            <n-text type="primary">
                                Try again later
                            </n-text>
                        </n-h4>
                    </n-space>
                </n-space>
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
    NSpace, NSkeleton,
    NIcon,
    NH1, NText, NH4
} from "naive-ui";
import { CloudOfflineOutline } from "@vicons/ionicons5";
import DashboardFooter from "@/components/core/footer.vue";

</script>

<style scoped>
.fullscreen {
    height: 99vh !important;
    width: 99vw !important;
}

.nav {
    width: 99vw !important;
}
</style>