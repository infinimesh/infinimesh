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
          :columns="columns"
          :data-source="namespaces"
          :loading="loading"
          rowKey="id"
          class="namespaces-table"
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
        </a-table>
      </a-col>
    </a-row>
  </div>
</template>

<script>
const columns = [
  {
    title: "Title",
    dataIndex: "name",
    sorter: true,
    scopedSlots: { customRender: "name" },
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
      loading: false,
    };
  },
  computed: {
    namespaces() {
      return this.$store.state.devices.namespaces;
    },
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
        placement: "bottomRight",
      });
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
table.namespaces-table {
  border-collapse: collapse;
}
.namespaces-table > table,
th,
td {
  border-bottom: 1px solid @primary-color !important;
  color: black;
}
</style>