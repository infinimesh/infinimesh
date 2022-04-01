<template>
  <n-layout>
    <n-layout-header>
      <dashboard-nav />
    </n-layout-header>
    <n-layout has-sider>
        <n-layout-sider content-style="padding: 24px;">
          Handian Bridge
        </n-layout-sider>
        <n-layout-content content-style="padding: 24px;">
          Pingshan Road
        </n-layout-content>
      </n-layout>
      <n-layout-footer position="absolute">Chengfu Road</n-layout-footer>
  </n-layout>
</template>

<script setup>
import { NLayout, NLayoutHeader, NLayoutContent, NLayoutSider, NLayoutFooter } from "naive-ui"
import DashboardNav from "@/components/dashboard/nav.vue"

import { inject } from "vue";
import { useAppStore } from "@/store/app";
const store = useAppStore()

const axios = inject('axios');

(async () => {
  axios.get('http://localhost:8000/accounts/me', {
    headers: {
      Authorization: `Bearer ${store.token}`
    }
  }).then(res => {
    store.me = res.data
  })
})()
</script>