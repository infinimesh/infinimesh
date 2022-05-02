<template>
  <n-spin :show="loading">
    <n-grid responsive="screen">
      <n-grid-item span="3">
        <n-h1 prefix="bar" align-text type="info">
          <n-text type="info"> Accounts </n-text>
        </n-h1>
      </n-grid-item>
      <n-grid-item span="15"> </n-grid-item>
      <n-grid-item span="3">
        <n-button strong secondary round type="info" @click="e => store.fetchAccounts()">
          <template #icon>
            <n-icon>
              <refresh-outline />
            </n-icon>
          </template>
          Refresh
        </n-button>
      </n-grid-item>
      <n-grid-item span="3">
        <account-create />
      </n-grid-item>
    </n-grid>
    <n-table :bordered="false" :single-line="true">
      <thead>
        <tr>
          <th>UUID</th>
          <th>Access</th>
          <th>Title</th>
          <th>Namespace</th>
          <th>Default Namespace</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="account in accounts" :key="account.uuid">
          <td>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button tertiary round :type="account.enabled ? 'success' : 'error'"
                  @click="handleCopyUUID(account.uuid)">
                  <template #icon>
                    <n-icon>
                      <copy-outline />
                    </n-icon>
                  </template>
                  {{ shortUUID(account.uuid) }}
                </n-button>
              </template>
              {{ account.uuid }}
            </n-tooltip>
          </td>
          <td>
            <AccessBadge :access="account.accessLevel" />
          </td>
          <td>
            <strong>
              {{ account.title }}
            </strong>
          </td>
          <td>
            {{ account.namespace }}
          </td>
          <td>
            {{ account.defaultNamespace || "-" }}
          </td>
          <td>
            <n-space>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button tertiary round :type="account.enabled ? 'error' : 'success'">
                    <template #icon>
                      <n-icon>
                        <ban-outline v-if="account.enabled" />
                        <checkmark-outline v-else />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                Click to {{ account.enabled ? "disable" : "enable" }} Account
              </n-tooltip>
              <n-popconfirm @positive-click="() => handleDelete(account.uuid)">
                <template #trigger>
                  <n-button v-if="account.accessLevel > 2" type="error" round secondary>Delete</n-button>
                </template>
                Are you sure about deleting this account?
              </n-popconfirm>
            </n-space>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-spin>
</template>

<script setup>
import { h, computed } from "vue";
import {
  NSpin,
  NTable,
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
import { CopyOutline, CheckmarkOutline, BanOutline, RefreshOutline } from "@vicons/ionicons5";
import { useAccountsStore } from "@/store/accounts";
import { storeToRefs } from "pinia";
import AccountCreate from "@/components/accounts/create-drawer.vue";

const store = useAccountsStore();
const { accounts_ns_filtered: accounts, loading } = storeToRefs(store);

const pool = computed(() => accounts);

store.fetchAccounts();

function shortUUID(uuid) {
  return uuid.substr(0, 8);
}

const accessLevels = {
  [0]: ["None", "error", undefined],
  [1]: ["Read", "error", undefined],
  [2]: ["Write", "warning", undefined],
  [3]: ["Admin", "success", undefined],
  [4]: ["Super-Admin", "success", "#8a2be2"],
};

const message = useMessage();
async function handleCopyUUID(uuid) {
  try {
    await navigator.clipboard.writeText(uuid);
    message.success("Account UUID copied to clipboard");
  } catch {
    message.error("Failed to copy Account UUID to clipboard");
  }
}

function AccessBadge(props, context) {
  let conf = accessLevels[props.access];
  return h(
    NButton,
    {
      secondary: true,
      round: true,
      type: conf[1],
      color: conf[2],
    },
    {
      default: () => conf[0],
    }
  );
}

const bar = useLoadingBar();
function handleDelete(uuid) {
  store.deleteAccount(uuid, bar)
}
}
</script>