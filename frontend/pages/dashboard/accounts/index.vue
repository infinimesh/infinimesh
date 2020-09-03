<template>
  <div id="accountsTable">
    <a-row>
      <a-col :span="23" :offset="1">
        <h1 class="lead">Accounts</h1>
      </a-col>
    </a-row>
    <a-row>
      <a-col :span="21" :offset="1">
        <a-table :columns="columns" :data-source="accounts" :loading="loading" rowKey="uid">
          <span slot="name" slot-scope="name">
            <b>{{ name }}</b>
          </span>
          <span slot="enabled" slot-scope="enabled">
            <a-row type="flex" justify="space-around">
              <a-icon type="bulb" :style="{color: enabled ? 'green' : 'red', fontSize: '24px' }" />
            </a-row>
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
    // width: "30%",
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
      this.$axios
        .get("/api/accounts")
        .then((res) => (this.accounts = res.data.accounts))
        .catch((e) => {
          if (e.response.status == 403) {
            this.$notification.error({
              message: "Oops",
              description: e.response.data.message,
            });
            this.$store.commit("window/noAccess", "dashboard-accounts");
            this.$router.push({ name: "dashboard-devices" });
          }
        });
    },
  },
};
</script>

<style scoped>
.ant-empty-description {
  color: lightgrey;
}
</style>