<template>
  <n-layout>
    <n-layout-header>
      <dashboard-nav />
    </n-layout-header>
    <n-layout has-sider :content-style="{ minHeight: '90vh' }">
      <n-layout-sider collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
        @mouseover="collapsed = false" @mouseleave="collapsed = true">
        <dashboard-menu :collapsed="collapsed" />
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
import { ref, computed } from "vue";
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NLayoutSider,
  NLayoutFooter,
} from "naive-ui";
import { useRoute } from "vue-router";
import { useAccountsStore } from "@/store/accounts";

import DashboardNav from "@/components/dashboard/nav.vue";
import DashboardMenu from "@/components/dashboard/menu.vue";
import DashboardFooter from "@/components/core/footer.vue";

const store = useAccountsStore();
const route = useRoute();

store.sync_me()

const collapsed = ref(true);
const noContentPadding = computed(() => {
  return route.name == 'DashboardMain'
})
</script>