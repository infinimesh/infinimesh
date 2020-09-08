<template>
  <div id="accountsTable">
    <a-row>
      <a-col :span="23" :offset="1">
        <h1 class="lead">Accounts</h1>
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
          <span slot="enabled" slot-scope="enabled">
            <a-row type="flex" justify="space-around">
              <a-icon type="bulb" :style="{color: enabled ? 'green' : 'red', fontSize: '24px' }" />
            </a-row>
          </span>
          <span slot="actions" slot-scope="text, account">
            <a-space>
              <a-button
                type="link"
                @click="toogleAccount(account)"
              >{{ account.enabled ? 'Disable' : 'Enable' }}</a-button>
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
const columns = [
  {
    title: "Username",
    dataIndex: "name",
    sorter: true,
    scopedSlots: { customRender: "name" },
  },
  {
    title: "Enabled",
    dataIndex: "enabled",
    sorter: true,
    width: "10%",
    scopedSlots: { customRender: "enabled" },
  },
  {
    title: "Actions",
    key: "actions",
    fixed: "right",
    width: "20%",
    scopedSlots: { customRender: "actions" },
  },
];

export default {
  data() {
    return {
      columns,
      accounts: [],
      loading: false,
    };
  },
  mounted() {
    this.getAccountsPool();
  },
  methods: {
    getAccountsPool() {
      const vm = this;
      vm.loading = true;
      vm.$axios
        .get("/api/accounts")
        .then((res) => (vm.accounts = res.data.accounts))
        .catch((e) => {
          if (e.response.status == 403) {
            vm.$notification.error({
              message: "Oops",
              description: e.response.data.message,
            });
            vm.$store.commit("window/noAccess", "dashboard-accounts");
            vm.$router.push({ name: "dashboard-devices" });
          }
        })
        .then(() => (vm.loading = false));
    },
    deleteAccount(account) {
      this.$notification.warning({
        message: "Coming soon",
        description: `Can't delete ${account.name}(${account.uid})`,
        placement: "bottomRight",
      });
    },
    toogleAccount(account) {
      const vm = this;
      vm.loading = true;
      vm.$axios({
        method: "patch",
        url: `/api/accounts/${account.uid}`,
        data: {
          enabled: !account.enabled,
        },
      })
        .then(() => {
          vm.$message.success(
            `Account successfuly ${account.enabled ? "disabled" : "enabled"}!`
          );
          vm.getAccountsPool();
        })
        .catch((e) => {
          vm.$notification.error({
            message: `Error ${
              account.enabled ? "disabling" : "enabling"
            } account`,
            description: e.response.data.message,
            placement: "bottomRight",
          });
        })
        .then(() => (vm.loading = false));
    },
  },
};
</script>

<style>
.ant-empty-description {
  color: lightgrey !important;
}
</style>
<style lang="less">
table.accounts-table {
  border-collapse: collapse;
}
.accounts-table > table,
th,
td {
  border-bottom: 1px solid @primary-color !important;
  color: black;
}
</style>