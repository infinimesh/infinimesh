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
        <n-button strong secondary round type="info" @click="e => store.fetchAccounts()">
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
            <strong>
              {{ account.title }}
            </strong>
          </td>
          <td>
            <AccessBadge :access="account.access.level" />
            <AccessBadge access="OWNER" v-if="account.access.role == 'OWNER'" left="5px" />
          </td>
          <td>
            {{ account.access.namespace }}
          </td>
          <td>
            {{ account.defaultNamespace || "-" }}
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

              <n-popconfirm @positive-click="() => handleDelete(account.uuid)">
                <template #trigger>
                  <n-button v-if="access_lvl_conv(account) > 2" type="error" round secondary>Delete</n-button>
                </template>
                <span>
                  Are you sure about deleting <b>{{ account.title }}'s</b> account?
                </span>
              </n-popconfirm>
            </n-space>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-spin>

  <set-credentials-modal :show="show_mc" @close="show_mc = false" :account="active_account" />
</template>

<script setup>
import { h, ref } from "vue";
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
import { CopyOutline, CheckmarkOutline, BanOutline, RefreshOutline, LockClosedOutline } from "@vicons/ionicons5";
import { useAccountsStore } from "@/store/accounts";
import { storeToRefs } from "pinia";
import AccountCreate from "@/components/accounts/create-drawer.vue";
import { access_lvl_conv } from "@/utils/access";
import setCredentialsModal from "@/components/accounts/set-credentials-modal.vue";

const store = useAccountsStore();
const { accounts_ns_filtered: accounts, loading } = storeToRefs(store);

store.fetchAccounts();

function shortUUID(uuid) {
  return uuid.substr(0, 8);
}

const accessLevels = {
  NONE: ["None", "error", undefined, "How did you get here??? Please, report this immideately"],
  READ: ["Read", "error", undefined, "You can only see this Account"],
  MGMT: ["Manage", "warning", undefined, "You can Manage this Account, for example enable/disable it"],
  ADMIN: ["Admin", "success", undefined, "You have the highest possible access to this Account"],
  ROOT: ["Super-Admin", "success", "#8a2be2", "You have the highest possible access to this Account"],
  OWNER: ["Owned", "success", "#8a2be2", "You are the owner of this Account, which gives you full access to it and right to delete it"]
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

function AccessBadge(props) {
  let conf = accessLevels[props.access];
  return h(
    NTooltip,
    {
      trigger: "hover",
      placement: "top",
    },
    {
      trigger: () => h(
        NButton,
        {
          secondary: true,
          round: true,
          type: conf[1],
          color: conf[2],
          style: {
            marginLeft: props.left
          }
        },
        {
          default: () => conf[0],
        }
      ),
      default: () => conf[3]
    }
  )
}

const bar = useLoadingBar();
function handleDelete(uuid) {
  store.deleteAccount(uuid, bar)
}
function handleToggleAccountEnabled(account) {
  store.toggle(account.uuid, bar);
}

const show_mc = ref(false);
const active_account = ref({})
</script>