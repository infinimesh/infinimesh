<template>
  <div>
    <a-table
      :loading="namespace.loading"
      :data-source="permissions"
      :columns="permissions_table_columns"
      :pagination="false"
      style="margin: 10px; width: 50%"
      :bordered="true"
      :locale="{ emptyText: 'No Permissions Found' }"
      :rowKey="(record, index) => `${record.account_id}-${index}`"
    >
      <span slot="account_name" slot-scope="account_name, record">
        <a-select
          v-if="record.editable"
          placeholder="Select Account"
          style="width: 100%"
          v-model="temp_permission.account_id"
          :options="
            accounts.map((acc) => {
              return { key: acc.uid, title: acc.name };
            })
          "
        >
        </a-select>
        <template v-else> {{ account_name }} </template>
      </span>
      <span slot="action" slot-scope="action, record">
        <a-select v-if="record.editable" v-model="temp_permission.action">
          <a-select-option
            :key="posAction[0]"
            v-for="posAction in Object.entries(actions)"
          >
            <a-tag :color="posAction[1]" slot="label">
              {{ posAction[0] }}
            </a-tag>
          </a-select-option>
          <a-select-option key="UNSELECTED" :disabled="true">
            <a-tag color="darkgrey" slot="label"> UNSELECTED </a-tag>
          </a-select-option>
        </a-select>
        <a-tag :color="actions[action]" v-else>
          {{ action }}
        </a-tag>
      </span>
      <span slot="actions" slot-scope="record">
        <a-button
          v-if="record.editable"
          type="link"
          @click="
            createPermission({ namespace: namespace.id, ...temp_permission })
          "
        >
          <a-icon type="save" style="font-size: 18px" />
        </a-button>
        <a-button v-else type="link" @click="deletePermission(record)">
          <a-icon type="delete" style="color: red; font-size: 18px" />
        </a-button>
      </span>
    </a-table>
    <a-row
      style="width: 50%"
      type="flex"
      justify="center"
      v-if="!temp_permission"
    >
      <a-col :span="16">
        <a-button
          type="primary"
          icon="plus"
          style="width: 100%"
          @click="addPermission"
          >Add Permission</a-button
        >
      </a-col>
    </a-row>
  </div>
</template>

<script>
import Vue from "vue";
import AccountControlMixin from "@/mixins/account-control";

const permissions_table_columns = [
  {
    title: "Account",
    dataIndex: "account_name",
    sorter: true,
    scopedSlots: { customRender: "account_name" },
  },
  {
    title: "Access",
    dataIndex: "action",
    width: "15%",
    scopedSlots: { customRender: "action" },
  },
  {
    title: "Actions",
    key: "actions",
    width: "15%",
    scopedSlots: { customRender: "actions" },
  },
];

export default Vue.component("namespace-permissions-table", {
  mixins: [AccountControlMixin],
  props: {
    namespace: { required: true },
  },
  data() {
    return {
      permissions_table_columns,
      accounts: [],
      temp_permission: null,
    };
  },
  computed: {
    permissions() {
      if (this.temp_permission) {
        return [...this.namespace.permissions, this.temp_permission];
      } else {
        return this.namespace.permissions;
      }
    },
  },
  created() {
    this.actions = {
      WRITE: "#eb2f96",
      READ: "#52c41a",
      NONE: "#5d8eb7",
    };
  },
  mounted() {
    this.getAccountsPool();
  },
  methods: {
    addPermission() {
      this.temp_permission = {
        action: "UNSELECTED",
        editable: true,
      };
    },
    createPermission({ namespace, account_id, action }) {
      let vm = this;
      vm.temp_permission = null;
      vm.$axios({
        method: "put",
        url: `/api/namespaces/${namespace}/permissions/${account_id}`,
        data: {
          action: action,
        },
      })
        .then(() => {
          vm.$notification.success({
            message: "Permission successfuly created.",
          });
          vm.$emit("refresh");
        })
        .catch((e) => {
          vm.$notification.error({
            message: "Failed to create permission",
            description: e.response.data.message,
          });
        })
        .then(() => {});
    },
    deletePermission({ namespace, account_id }) {
      let vm = this;
      vm.$axios({
        method: "delete",
        url: `/api/namespaces/${namespace}/permissions/${account_id}`,
      })
        .then(() => {
          vm.$notification.success({
            message: "Permission successfuly deleted.",
          });
          vm.$emit("refresh");
        })
        .catch((e) => {
          vm.$notification.error({
            message: "Failed to delete permission",
            description: e.response.data.message,
          });
        })
        .then(() => {});
    },
  },
});
</script>