<template>
  <n-spin :show="loading">
    <n-grid item-responsive>
      <n-grid-item span="24 500:12 1000:4">
        <n-h1 prefix="bar" align-text type="info">
          <n-text type="info"> Accounts </n-text>
        </n-h1>
      </n-grid-item>
      <n-grid-item span="0 600:2 700:4 1000:12 1400:14"> </n-grid-item>
      <n-grid-item span="12 500:6 600:5 700:4 1000:4 1400:2">
        <n-button strong secondary round type="info" @click="e => store.fetchAccounts(true)">
          <template #icon>
            <n-icon>
              <refresh-outline />
            </n-icon>
          </template>
          Refresh
        </n-button>
      </n-grid-item>
      <n-grid-item span="12 500:6 600:5 700:4 1000:4 1400:2">
        <account-create />
      </n-grid-item>
    </n-grid>
    <n-table :bordered="false" :single-line="true" style="margin-top: 10px">
      <thead>
        <tr>
          <th>UUID</th>
          <th>Title</th>
          <th>Access</th>
          <th>Namespace</th>
          <th>Default Namespace</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="account in accounts" :key="account.uuid">
          <td>
            <uuid-badge :uuid="account.uuid" :type="account.enabled ? 'success' : 'error'" />
          </td>
          <td>
            <strong>
              {{ account.title }}
            </strong>
          </td>
          <td>
            <access-badge :access="account.access.level" />
            <access-badge access="OWNER" v-if="account.access.role == 'OWNER'" left="5px" />
          </td>
          <td>
            {{ nss.namespaces[account.access.namespace]?.title || account.access.namespace }}
          </td>
          <td>
            {{ nss.namespaces[account.defaultNamespace]?.title || "-" }}
          </td>
          <td>
            <n-space>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button tertiary circle :type="account.enabled ? 'error' : 'success'"
                    @click="e => handleToggleAccountEnabled(account)">
                    <template #icon>
                      <n-icon>
                        <ban-outline v-if="account.enabled" />
                        <checkmark-outline v-else />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                <span>Click to {{ account.enabled ? "disable" : "enable" }} <b>{{ account.title }}'s</b> Account</span>
              </n-tooltip>

              <n-tooltip v-if="access_lvl_conv(account) > 3 || account.access.role == 'OWNER'" trigger="hover">
                <template #trigger>
                  <n-button type="warning" @click="() => { active_account = account; show_mc = true }" tertiary circle>
                    <template #icon>
                      <n-icon>
                        <lock-closed-outline />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                <span>Click to manage <b>{{ account.title }}'s</b> credentials</span>
              </n-tooltip>

              <n-button round secondary type="success" @click="e => handleLoginAs(account)">
                Login as
              </n-button>
              <acc-delete :o="account" :deletables="() => showDeletables(account.uuid)"
                @confirm="() => handleDelete(account.uuid)" type="account" />
            </n-space>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-spin>

  <set-credentials-modal :show="show_mc" @close="show_mc = false" :account="active_account" />
</template>

<script setup>
import { ref } from "vue";
import {
  NSpin,
  NTable,
  NButton,
  NIcon,
  NSpace,
  NTooltip,
  NGrid,
  NGridItem,
  NH1,
  NText,
  useLoadingBar, useMessage
} from "naive-ui";
import { CheckmarkOutline, BanOutline, RefreshOutline, LockClosedOutline } from "@vicons/ionicons5";
import { useAppStore } from "@/store/app";
import { useAccountsStore } from "@/store/accounts";
import { useNSStore } from "@/store/namespaces";
import { storeToRefs } from "pinia";

import AccountCreate from "@/components/accounts/create-drawer.vue";

import UuidBadge from "@/components/core/uuid-badge.vue";

import AccessBadge from "@/components/core/access-badge"
import { access_lvl_conv } from "@/utils/access";
import setCredentialsModal from "@/components/accounts/set-credentials-modal.vue";
import AccDelete from "@/components/core/recursive-delete-modal.vue";


const store = useAccountsStore();
const { accounts_ns_filtered: accounts, loading } = storeToRefs(store);

store.fetchAccounts(true);

const show_mc = ref(false);
const active_account = ref({})

const nss = useNSStore()

async function showDeletables(uuid) {
  return store.deletables(uuid)
}

const bar = useLoadingBar();
function handleToggleAccountEnabled(account) {
  store.toggle(account.uuid, bar);
}

const message = useMessage()
async function handleDelete(uuid) {
  loading.value = true
  try {
    await store.deleteAccount(uuid, bar)
    delete store.accounts[uuid]
    message.success("Account successfuly deleted")
  } catch (e) {
    message.error("Failed to delete account: " + e.response.statusText)
  }
}

const as = useAppStore()
async function handleLoginAs(account) {
  try {
    const { data } = await store.tokenFor(account.uuid)
    const params = { token: data.token, title: account.title, back_token: as.token }

    window.open(window.location.origin + '/#login?a=' + btoa(JSON.stringify(params)), '_blank', { incognito: true })
  } catch (e) {
    message.error("Failed to get token: " + e.response.statusText)
  }
}
</script>