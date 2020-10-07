<template>
  <div id="namespacesTable">
    <a-row>
      <a-col :span="23" :offset="1">
        <h1 class="lead">Namespaces</h1>
      </a-col>
    </a-row>
    <a-row>
      <a-col :span="21" :offset="1">
        <a-table
          :columns="namespaces_table_columns"
          :data-source="namespaces"
          :loading="loading"
          rowKey="id"
          class="namespaces-table"
          :expandRowByClick="true"
          @expand="loadNamespacePermissions"
        >
          <span slot="name" slot-scope="name">
            <b>{{ name }}</b>
          </span>
          <span slot="actions" slot-scope="text, namespace">
            <a-space>
              <a-button type="link" @click="deleteNamespace(namespace)">
                <a-icon type="delete" style="color: red; font-size: 18px" />
              </a-button>
            </a-space>
          </span>

          <a-table
            slot="expandedRowRender"
            slot-scope="record"
            :loading="record.loading"
            :data-source="record.permissions"
            :columns="permissions_table_columns"
            :pagination="false"
            style="margin: 10px; width: 50%"
            :bordered="true"
            :rowKey="(record, index) => `${record.account_id}-${index}`"
          >
            <span slot="action" slot-scope="action">
              <a-tag :color="actionColors[action]">
                {{ action }}
              </a-tag>
            </span>
          </a-table>
        </a-table>
      </a-col>
    </a-row>
  </div>
</template>

<script>
const namespaces_table_columns = [
  {
    title: "Title",
    dataIndex: "name",
    sorter: true,
    scopedSlots: { customRender: "name" },
  },
  {
    title: "Actions",
    key: "actions",
    width: "10%",
    scopedSlots: { customRender: "actions" },
  },
];
const permissions_table_columns = [
  {
    title: "Account",
    dataIndex: "account_name",
    sorter: true,
  },
  {
    title: "Access",
    dataIndex: "action",
    width: "10%",
    scopedSlots: { customRender: "action" },
  },
];

export default {
  data() {
    return {
      namespaces_table_columns,
      permissions_table_columns,
      loading: false,
    };
  },
  computed: {
    namespaces() {
      return this.$store.state.devices.namespaces;
    },
  },
  created() {
    this.actionColors = {
      WRITE: "#eb2f96",
      READ: "#52c41a",
      NONE: "#5d8eb7",
    };
  },
  mounted() {
    this.getNamespacesPool();
  },
  methods: {
    async getNamespacesPool() {
      this.loading = true;
      await this.$store.dispatch("devices/getNamespaces");
      this.loading = false;
    },
    deleteNamespace(namespace) {
      this.$notification.warning({
        message: "Coming soon",
        description: `Can't delete ${namespace.name}(${namespace.id})`,
      });
    },
    loadNamespacePermissions(expanded, ns) {
      console.log(expanded, ns);
      if (expanded) {
        this.$store.dispatch("devices/getNamespacePermissions", ns);
      }
    },
  },
};
</script>

<style>
.ant-empty-description {
  color: lightgrey !important;
}
table.namespaces-table {
  border-collapse: collapse;
}
.namespaces-table > table,
th,
td {
  border-bottom: 1px solid var(--primary-color) !important;
  color: black;
}
</style>