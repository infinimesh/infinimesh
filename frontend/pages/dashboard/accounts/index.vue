<template>
  <div id="accountsTable">
    <a-row type="flex" align="middle">
      <a-col :span="12" :offset="1">
        <h1 class="lead">Accounts</h1>
      </a-col>
      <a-col :span="3" :offset="6">
        <a-row type="flex" justify="end">
          <a-button
            type="primary"
            icon="plus"
            @click="createAccountDrawerVisible = true"
            >Create Account</a-button
          >
        </a-row>
        <account-add
          :active="createAccountDrawerVisible"
          @cancel="createAccountDrawerVisible = false"
          @add="handleAccountAdd"
        />
      </a-col>
    </a-row>
    <a-row>
      <a-col :span="21" :offset="1">
        <a-table
          :columns="columns"
          :data-source="accounts"
          :loading="loading"
          rowKey="uid"
          class="accounts-table"
        >
          <span slot="name" slot-scope="name">
            <b>{{ name }}</b>
          </span>
          <span slot="is_admin" slot-scope="is_admin">
            <a-row type="flex" justify="space-around">
              <a-icon
                :type="is_admin ? 'check-circle' : 'close-circle'"
                :style="{ color: is_admin ? 'green' : 'red', fontSize: '24px' }"
              />
            </a-row>
          </span>
          <span slot="enabled" slot-scope="enabled">
            <a-row type="flex" justify="space-around">
              <a-icon
                type="bulb"
                :style="{ color: enabled ? 'green' : 'red', fontSize: '24px' }"
              />
            </a-row>
          </span>
          <span slot="actions" slot-scope="text, account">
            <a-space>
              <a-button type="link" @click="resetAccountPassword(account)"
                >Reset password</a-button
              >
              <account-reset-password
                v-if="selectedAccount"
                :active="resetAccountPasswordVisible"
                :account="selectedAccount"
                @cancel="resetAccountPasswordVisible = false"
                @reset="handleResetAccountPassword"
              />

              <a-divider type="vertical" />

              <a-button type="link" @click="toogleAccount(account)">{{
                account.enabled ? "Disable" : "Enable"
              }}</a-button>

              <a-divider type="vertical" />

              <a-button type="link" @click="deleteAccount(account)">
                <a-icon type="delete" style="color: red; font-size: 18px" />
              </a-button>
            </a-space>
          </span>
        </a-table>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import AccountAdd from "@/components/account/Add.vue";
import AccountResetPassword from "@/components/account/ResetPassword.vue";

import AccountControlMixin from "@/mixins/account-control";

const columns = [
  {
    title: "Username",
    dataIndex: "name",
    sorter: true,
    scopedSlots: { customRender: "name" },
  },
  {
    title: "Admin",
    dataIndex: "is_admin",
    sorter: true,
    width: "7%",
    scopedSlots: { customRender: "is_admin" },
  },
  {
    title: "Enabled",
    dataIndex: "enabled",
    sorter: true,
    width: "7%",
    scopedSlots: { customRender: "enabled" },
  },
  {
    title: "Actions",
    key: "actions",
    fixed: "right",
    width: "15%",
    scopedSlots: { customRender: "actions" },
  },
];

export default {
  mixins: [AccountControlMixin],
  components: {
    AccountAdd,
    AccountResetPassword,
  },
  computed: {
    user() {
      return this.$store.getters.loggedInUser;
    },
  },
  data() {
    return {
      columns,
      accounts: [],
      loading: false,

      createAccountDrawerVisible: false,

      resetAccountPasswordVisible: false,
      selectedAccount: null,
    };
  },
  mounted() {
    this.getAccountsPool();
  },
  methods: {
    resetAccountPassword(account) {
      this.selectedAccount = account;
      this.resetAccountPasswordVisible = true;
    },
    handleResetAccountPassword(password) {
      this.resetAccountPasswordVisible = false;
      this.updateAccount(
        this.selectedAccount.uid,
        { password: password },
        "Password changed successfuly",
        "Reset password failed"
      );
    },
  },
};
</script>

<style>
.ant-empty-description {
  color: lightgrey !important;
}
table.accounts-table {
  border-collapse: collapse;
}
.accounts-table > table,
th,
td {
  border-bottom: 1px solid var(--primary-color) !important;
  color: black;
}
</style>
