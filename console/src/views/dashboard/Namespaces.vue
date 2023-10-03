<template>
  <n-spin :show="loading">
    <n-grid item-responsive>
      <n-grid-item span="24 400:12 550:8 1000:4">
        <n-h1 prefix="bar" align-text type="info">
          <n-text type="info"> Namespaces </n-text>
        </n-h1>
      </n-grid-item>
      <n-grid-item span="0 1000:10 1400:12"> </n-grid-item>
      <n-grid-item span="8 400:12 550:6 1000:3 1400:2" align="end">
        <n-button strong secondary round type="info" @click="refresh">
          <template #icon>
            <n-icon>
              <refresh-outline />
            </n-icon>
          </template>
          Refresh
        </n-button>
      </n-grid-item>
      <n-grid-item span="16 400:24 550:10 1000:7 1400:6" align="end">
        <ns-create />
      </n-grid-item>
    </n-grid>
    <n-table :bordered="false" :single-line="true" style="margin-top: 10px">
      <n-thead>
        <n-tr>
          <th></th>
          <th>UUID</th>
          <th>Title</th>
          <th>Access</th>
          <th>Actions</th>
        </n-tr>
      </n-thead>
      <tbody>
        <template v-for="ns in pool.admin" :key="ns.uuid">
          <n-tr @click="expand.has(ns.uuid) ? expand.delete(ns.uuid) : expand.add(ns.uuid)">
            <td>
              <n-icon>
                <chevron-down-outline v-if="expand.has(ns.uuid)" />
                <chevron-forward-outline v-else />
              </n-icon>
            </td>
            <td>
              <uuid-badge :uuid="ns.uuid" />
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
                <config-edit-modal :o="ns" @submit="handleConfigUpdate"
                  v-if="access_lvl_conv(ns) > 3 || ns.access.role == 'OWNER'" />

                <n-button type="success" round secondary @click.stop.prevent="setNSAndGo(ns.uuid, 'Devices')">
                  Devices
                </n-button>
                <n-button type="info" round secondary @click.stop.prevent="setNSAndGo(ns.uuid, 'Accounts')">
                  Accounts
                </n-button>
                <ns-delete :o="ns" :deletables="() => showDeletables(ns.uuid)" @confirm="() => handleDelete(ns.uuid)" />
              </n-space>
            </td>
          </n-tr>
          <ns-joins v-if="expand.has(ns.uuid)" :namespace="ns.uuid"
            :admin="ns.access.role == 'OWNER' || ns.access.level == 'ROOT'" />
        </template>
        <n-tr v-if="pool.user && pool.user.length">
          <td colspan="5" align="center">
            <span>
              Namespaces below are those you don't have admin access to
            </span>
          </td>
        </n-tr>
        <n-tr v-for="ns in pool.user" :key="ns.uuid">
          <td></td>
          <td>
            <uuid-badge :uuid="ns.uuid" type="default" />
          </td>
          <td>
            <strong>
              {{ ns.title }}
            </strong>
          </td>
          <td>
            <access-badge :access="ns.access.level" namespace />
          </td>
          <td>
            <n-space>
              <n-button type="success" round secondary @click.stop.prevent="setNSAndGo(ns.uuid, 'Devices')">
                Devices
              </n-button>
              <n-button type="info" round secondary @click.stop.prevent="setNSAndGo(ns.uuid, 'Accounts')">
                Accounts
              </n-button>
            </n-space>
          </td>
        </n-tr>
      </tbody>
    </n-table>
  </n-spin>
</template>

<script setup>
import { ref, computed, defineAsyncComponent } from "vue"
import { useRouter } from "vue-router"
import {
  NSpin, NTable, NThead, NTr,
  NButton, NIcon, NSpace, NGrid,
  NGridItem, NH1, NText, useMessage
} from "naive-ui";

import { useNSStore } from "@/store/namespaces";

import { storeToRefs } from "pinia";
import { access_lvl_conv } from "@/utils/access";
import { groupBy } from "lodash";

const RefreshOutline = defineAsyncComponent(() => import("@vicons/ionicons5/RefreshOutline"))
const ChevronForwardOutline = defineAsyncComponent(() => import("@vicons/ionicons5/ChevronForwardOutline"))
const ChevronDownOutline = defineAsyncComponent(() => import("@vicons/ionicons5/ChevronDownOutline"))

const UuidBadge = defineAsyncComponent(() => import("@/components/core/uuid-badge.vue"))
const AccessBadge = defineAsyncComponent(() => import("@/components/core/access-badge"))

const NsCreate = defineAsyncComponent(() => import("@/components/namespaces/create-action-button.vue"))
const NsJoins = defineAsyncComponent(() => import("@/components/namespaces/joins.vue"))
const NsDelete = defineAsyncComponent(() => import("@/components/core/recursive-delete-modal.vue"))

const ConfigEditModal = defineAsyncComponent(() => import("@/components/core/config-edit-modal.vue"))

const store = useNSStore();
const { loading, namespaces_list: namespaces } = storeToRefs(store);

const pool = computed(() => groupBy(namespaces.value, (e) => {
  if ((e.access ?? { role: "" }).role == "OWNER" || access_lvl_conv(e) >= 3) {
    return "admin"
  }
  return "user"
}))

const expand = ref(new Set())

const router = useRouter()
function setNSAndGo(ns, route) {
  store.selected = ns
  router.push({ name: route })
}

function refresh() {
  expand.value = new Set()
  store.fetchNamespaces(true)
}

refresh()

async function showDeletables(uuid) {
  return store.deletables(uuid)
}

const message = useMessage()
async function handleDelete(uuid) {
  loading.value = true
  try {
    await store.delete(uuid)
    message.success("Namespace successfuly deleted")
  } catch (e) {
    message.error("Failed to delete namespace: " + e.response.statusText)
  }
  refresh()
}

async function handleConfigUpdate(ns) {
  try {
    await store.update(ns)
    message.success("Namespace config updated")
  } catch (e) {
    message.error("Failed to update namespace config: " + e.response.statusText)
  }
  refresh()
}

</script>