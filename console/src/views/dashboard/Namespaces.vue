<template>
  <n-spin :show="loading">
    <n-grid item-responsive>
      <n-grid-item span="24 500:12 1000:4">
        <n-h1 prefix="bar" align-text type="info">
          <n-text type="info"> Namespaces </n-text>
        </n-h1>
      </n-grid-item>
      <n-grid-item span="0 600:2 700:4 1000:12 1400:14"> </n-grid-item>
      <n-grid-item span="12 500:6 600:5 700:4 1000:4 1400:2">
        <n-button strong secondary round type="info" @click="e => store.fetchNamespaces()">
          <template #icon>
            <n-icon>
              <refresh-outline />
            </n-icon>
          </template>
          Refresh
        </n-button>
      </n-grid-item>
      <n-grid-item span="12 500:6 600:5 700:4 1000:4 1400:2">
        <!-- <ns-create /> -->
      </n-grid-item>
    </n-grid>
    <n-table :bordered="false" :single-line="true" style="margin-top: 10px">
      <n-thead>
        <n-tr>
          <th>UUID</th>
          <th>Title</th>
          <th>Access</th>
          <th>Actions</th>
        </n-tr>
      </n-thead>
      <tbody>
        <n-tr v-for="ns in pool.admin" :key="ns.uuid">
          <td>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button tertiary round type="info" @click="handleCopyUUID(ns.uuid)">
                  <template #icon>
                    <n-icon>
                      <copy-outline />
                    </n-icon>
                  </template>
                  {{ shortUUID(ns.uuid) }}
                </n-button>
              </template>
              {{ ns.uuid }}
            </n-tooltip>
          </td>
          <td>
            <strong>
              {{ ns.title }}
            </strong>
          </td>
          <td>
            <access-badge :access="ns.access.level" namespace />
            <access-badge access="OWNER" v-if="ns.access.role == 'OWNER'" left="5px" namespace />
          </td>
          <td>
            <n-space>
              <n-popconfirm @positive-click="() => handleDelete(ns.uuid)">
                <template #trigger>
                  <n-button v-if="ns.access.role == 'OWNER'" type="error" round secondary>Delete</n-button>
                </template>
                <span>
                  Are you sure about deleting <b>{{ ns.title }}'s</b> ns?
                </span>
              </n-popconfirm>
            </n-space>
          </td>
        </n-tr>
        <n-tr>
          <td colspan="4" align="center">
            <span>
              Namespaces below are those you don't have admin access to
            </span>
          </td>
        </n-tr>
        <n-tr v-for="ns in pool.user" :key="ns.uuid">
          <td>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button tertiary round type="default" @click="handleCopyUUID(ns.uuid)">
                  <template #icon>
                    <n-icon>
                      <copy-outline />
                    </n-icon>
                  </template>
                  {{ shortUUID(ns.uuid) }}
                </n-button>
              </template>
              {{ ns.uuid }}
            </n-tooltip>
          </td>
          <td>
            <strong>
              {{ ns.title }}
            </strong>
          </td>
          <td>
            <access-badge :access="ns.access.level" namespace />
          </td>
          <td></td>
        </n-tr>
      </tbody>
    </n-table>
  </n-spin>
</template>

<script setup>
import { computed } from "vue"
import {
  NSpin,
  NTable,
  NThead,
  NTr,
  NButton,
  NIcon,
  NSpace,
  useMessage,
  NTooltip,
  NPopconfirm,
  NGrid,
  NGridItem,
  NH1,
  NText, useLoadingBar
} from "naive-ui";
import { CopyOutline, CheckmarkOutline, BanOutline, RefreshOutline, LockClosedOutline } from "@vicons/ionicons5";
import { useNSStore } from "@/store/namespaces";
import { storeToRefs } from "pinia";
import { access_lvl_conv } from "@/utils/access";
import { groupBy } from "lodash"

import AccessBadge from "@/components/core/access-badge"

function shortUUID(uuid) {
  return uuid.substr(0, 8);
}

const store = useNSStore();
const { loading, namespaces } = storeToRefs(store);

const pool = computed(() => groupBy(namespaces.value, (e) => {
  if (e.access.role == "OWNER" || access_lvl_conv(e) >= 3) {
    return "admin"
  }
  return "user"
}))

store.fetchNamespaces()
</script>