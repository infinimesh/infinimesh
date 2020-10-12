<template>
  <a-drawer
    title="Create New Namespace"
    :visible="active"
    :width="drawerSize"
    :headerStyle="{ fontSize: '4rem' }"
    @close="$emit('cancel')"
  >
    <a-form-model
      :model="namespace"
      :rules="rules"
      :label-col="{ xs: 24, sm: 6, md: 6, lg: 6 }"
      :wrapper-col="{ xs: 24, sm: 16, md: 18, lg: { span: 14, offset: 1 } }"
      ref="namespaceAddForm"
    >
      <a-form-model-item prop="name" label="Name">
        <a-input v-model="namespace.name" />
      </a-form-model-item>
    </a-form-model>
    <div
      :style="{
        position: 'absolute',
        right: 0,
        bottom: 0,
        width: '100%',
        borderTop: '1px solid #e9e9e9',
        padding: '10px 16px',
        textAlign: 'right',
        zIndex: 1,
      }"
      class="add-drawer-actions-row"
    >
      <a-button
        :style="{ marginRight: '8px' }"
        @click="setDefault"
        class="ant-btn-danger"
        >Reset</a-button
      >
      <a-button
        :style="{ marginRight: '8px' }"
        @click="$emit('cancel')"
        type="primary"
        >Cancel</a-button
      >
      <a-button type="success" @click="handleSubmit">Submit</a-button>
    </div>
  </a-drawer>
</template>

<script>
import Vue from "vue";

import drawerSizeMixin from "@/mixins/drawer-size.js";

export default Vue.component("namespace-add", {
  mixins: [drawerSizeMixin],
  props: {
    active: {
      required: true,
    },
  },
  data() {
    return {
      namespace: {},
      rules: {
        name: [
          {
            required: true,
            message: "Please, fill in the namespace name",
          },
          {
            pattern: /^[a-zA-Z0-9\-_]*$/,
            message:
              "Namespace name can contain only alphanumeric characters, hyphens and underscores",
          },
        ],
      },
    };
  },
  watch: {
    active: "setDefault",
  },
  mounted() {
    this.setDefault();
  },
  methods: {
    setDefault() {
      this.namespace = {
        name: "",
      };
    },
    handleSubmit() {
      let form = this.$refs["namespaceAddForm"];
      let errors = [];
      form.validateField(Object.keys(this.namespace), (err) => {
        if (err) {
          errors.push(err);
        }
      });
      if (errors.length === 0) {
        this.$emit("add", this.namespace);
      }
    },
  },
});
</script>