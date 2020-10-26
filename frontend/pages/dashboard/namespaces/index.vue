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
          @expand="loadNamespacePermissions"
        >
          <span slot="name" slot-scope="name, namespace">
            <a-input
              v-if="namespace.editable"
              style="width: 50%"
              :default-value="namespace.name"
              @change="$store.commit('devices/update_namespace', namespace)"
              placeholder="Enter new name"
            />
            <b v-else>{{ name }}</b>
          </span>
          <span slot="actions" slot-scope="text, namespace">
            <a-space>
              <template v-if="namespace.editable">
                <a-button type="link" @click="renameNamespace(namespace)">
                  <a-icon type="save" style="font-size: 18px" />
                </a-button>
                <a-button
                  type="link"
                  v-if="namespace.editable"
                  @click="getNamespacesPool"
                >
                  <a-icon type="close" style="color: red; font-size: 18px" />
                </a-button>
              </template>
              <template v-else>
                <a-button
                  type="link"
                  style="font-size: 18px"
                  @click="
                    $store.commit('devices/update_namespace', {
                      ...namespace,
                      editable: true,
                    })
                  "
                >
                  <a-icon type="edit" />
                </a-button>

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
                <a-button
                  v-else
                  type="link"
                  @click="deleteNamespace(namespace)"
                >
                  <a-icon type="delete" style="color: red; font-size: 18px" />
                </a-button>
              </template>
            </a-space>
          </span>

          <span slot="expandedRowRender" slot-scope="record">
            <namespace-permissions-table
              :namespace="record"
              @refresh="loadNamespacePermissions(true, record)"
            />
          </span>
        </a-table>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import NamespaceAdd from "@/components/namespace/Add";
import NamespacePermissionsTable from "@/components/namespace/PermissionsTable";

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

export default {
  components: {
    NamespaceAdd,
    NamespacePermissionsTable,
  },
  data() {
    return {
      namespaces_table_columns,
      loading: false,

      createNamespaceDrawerVisible: false,
    };
  },
  computed: {
    namespaces() {
      return this.$store.state.devices.namespaces;
    },
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
    updateNamespace(ns_id, patch, { success, error, always }) {
      this.$axios({
        url: `/api/namespaces/${ns_id}`,
        method: "patch",
        data: patch,
      })
        .then((res) => {
          if (success) success(res);
        })
        .catch((e) => {
          if (error) error(e);
        })
        .then(() => {
          if (always) always();
        });
    },
    renameNamespace(ns) {
      const vm = this;
      vm.updateNamespace(
        ns.id,
        { name: ns.name },
        {
          success: () => {
            vm.$message.success("Namespace successfuly renamed!");
          },
          error: (e) => {
            vm.$notification.error({
              message: "Error renaming namespace " + namespace.name,
              description: e.response.data.message,
            });
          },
          always: () => {
            vm.getNamespacesPool();
          },
        }
      );
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
      this.updateNamespace(
        namespace.id,
        {
          markfordeletion: false,
        },
        {
          success: () => {
            vm.$message.success("Namespace successfuly restored!");
          },
          error: (e) => {
            vm.$notification.error({
              message: "Error restoring namespace " + namespace.name,
              description: e.response.data.message,
            });
          },
          always: () => {
            vm.getNamespacesPool();
          },
        }
      );
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