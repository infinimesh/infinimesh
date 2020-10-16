<template>
  <div id="namespacesTable">
    <a-row type="flex" align="middle">
      <a-col :span="12" :offset="1">
        <h1 class="lead">Namespaces</h1>
      </a-col>
      <a-col :span="3" :offset="6">
        <a-row type="flex" justify="end">
          <a-button
            type="primary"
            icon="plus"
            @click="createNamespaceDrawerVisible = true"
            >Create Namespace</a-button
          >
        </a-row>
        <namespace-add
          :active="createNamespaceDrawerVisible"
          @cancel="createNamespaceDrawerVisible = false"
          @add="handleNamespaceAdd"
        />
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
              <a-tooltip
                v-if="namespace.markfordeletion"
                :title="`Going to be deleted ${deletionTime(namespace)}`"
                placement="left"
              >
                <a-button type="link" @click="restoreNamespace(namespace)">
                  <a-icon
                    type="redo"
                    style="color: var(--switch-color); font-size: 18px"
                  />
                  Restore
                </a-button>
              </a-tooltip>

              <a-button v-else type="link" @click="deleteNamespace(namespace)">
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
            :locale="{ emptyText: 'No Permissions Found' }"
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
import NamespaceAdd from "@/components/namespace/Add";

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
    width: "15%",
    scopedSlots: { customRender: "action" },
  },
];

export default {
  components: {
    NamespaceAdd,
  },
  data() {
    return {
      namespaces_table_columns,
      permissions_table_columns,
      loading: false,

      createNamespaceDrawerVisible: false,
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
    if (this.$route.query.create) this.createNamespaceDrawerVisible = true;
  },
  methods: {
    async getNamespacesPool() {
      this.loading = true;
      await this.$store.dispatch("devices/getNamespaces");
      this.loading = false;
    },
    deleteNamespace(namespace) {
      const vm = this;
      vm.$axios({
        url: `/api/namespaces/${namespace.id}/false`,
        method: "delete",
      })
        .then(() => {
          vm.$message.success("Namespace successfuly deleted!");
          vm.getNamespacesPool();
        })
        .catch((e) => {
          vm.$notification.error({
            message: "Error deleting namespace " + namespace.name,
            description: e.response.data.message,
          });
        });
    },
    restoreNamespace(namespace) {
      const vm = this;
      vm.$axios({
        url: `/api/namespaces/${namespace.id}`,
        method: "patch",
        data: {
          markfordeletion: false,
        },
      })
        .then(() => {
          vm.$message.success("Namespace successfuly restored!");
          vm.getNamespacesPool();
        })
        .catch((e) => {
          vm.$notification.error({
            message: "Error restoring namespace " + namespace.name,
            description: e.response.data.message,
          });
        });
    },
    loadNamespacePermissions(expanded, ns) {
      if (expanded) {
        this.$store.dispatch("devices/getNamespacePermissions", ns);
      }
    },
    handleNamespaceAdd(namespace) {
      const vm = this;
      vm.$axios({
        method: "post",
        url: "/api/namespaces",
        data: namespace,
      })
        .then(() => {
          vm.$notification.success({
            message: "Namespace created successfuly",
          });
          vm.createNamespaceDrawerVisible = false;
          vm.getNamespacesPool();
        })
        .catch((err) => {
          this.$notification.error({
            message: "Failed to create a namespace",
            description: `Response: ${err.response.data.message}`,
            duration: 10,
          });
        });
    },
    deletionTime(namespace) {
      let delete_init_date = new Date(namespace.deleteinitiationtime);
      delete_init_date.setDate(delete_init_date.getDate() + 14);
      return "on " + delete_init_date;
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