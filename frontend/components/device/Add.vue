<template>
  <a-drawer
    title="Create New Device"
    :visible="active"
    :width="drawerSize"
    :headerStyle="{ fontSize: '4rem' }"
    @close="$emit('cancel')"
  >
    <a-form-model
      :model="device"
      :rules="rules"
      :label-col="{ xs: 24, sm: 6, md: 6, lg: 6 }"
      :wrapper-col="{ xs: 24, sm: 16, md: 18, lg: { span: 14, offset: 1 } }"
      ref="deviceAddForm"
    >
      <a-form-model-item prop="name" label="Name">
        <a-input v-model="device.name" />
      </a-form-model-item>
      <a-form-model-item label="Namespace" prop="namespace">
        <a-select
          v-model="device.namespace"
          placeholder="please select device namespace"
        >
          <a-select-option
            :value="ns.id"
            :key="ns.id"
            v-for="ns in namespaces"
            >{{ ns.name }}</a-select-option
          >
        </a-select>
      </a-form-model-item>
      <a-form-model-item prop="tags" label="Tags">
        <a-select
          mode="tags"
          :token-separators="[',']"
          v-model="device.tags"
          placeholder="Enter a comma-separated list of tags, e.g. tag1, tag2"
        />
      </a-form-model-item>
      <a-form-model-item label="Enabled" prop="enabled">
        <a-switch v-model="device.enabled" />
      </a-form-model-item>
      <a-form-model-item prop="certificate" label="Certificate">
        <a-tabs
          default-active-key="paste"
          id="certificate-tabs"
          v-model="certificate_tab"
        >
          <a-tab-pane key="paste">
            <span slot="tab"> <a-icon type="copy" />Paste </span>
            <a-textarea
              v-model="device.certificate.pem_data"
              placeholder="Paste your certificate"
              :autoSize="{ minRows: 10, maxRows: 30 }"
            />
          </a-tab-pane>
          <a-tab-pane key="upload">
            <span slot="tab"> <a-icon type="upload" />Upload </span>
            <a-upload-dragger
              name="certificate"
              accept="pem, crt, pub"
              :beforeUpload="handleUploadCertificate"
              :showUploadList="false"
            >
              <p class="ant-upload-drag-icon">
                <a-icon type="inbox" />
              </p>
              <p class="ant-upload-text">
                Click or drag file to this area to upload
              </p>
            </a-upload-dragger>
          </a-tab-pane>
        </a-tabs>
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
      id="deviceAddDrawerActionsRow"
    >
      <a-button
        :style="{ marginRight: '8px' }"
        @click="setDefault"
        class="ant-btn-danger"
        >Reset</a-button
      >
      <a-button :style="{ marginRight: '8px' }" @click="$emit('cancel')"
        >Cancel</a-button
      >
      <a-button type="success" @click="handleSubmit">Submit</a-button>
    </div>
  </a-drawer>
</template>

<script>
import Vue from "vue";

import drawerSizeMixin from "@/mixins/drawer-size.js";

export default Vue.component("device-add", {
  mixins: [drawerSizeMixin],
  props: {
    active: {
      required: true,
    },
  },
  computed: {
    namespaces() {
      return this.$store.state.devices.namespaces;
    },
  },
  watch: {
    active: "setDefault",
  },
  data() {
    return {
      device: {
        name: "",
        tags: [],
        namespace: "",
        enabled: false,
        certificate: { pem_data: "" },
      },
      certificate_tab: "upload",
      rules: {
        name: [
          { required: true, message: "Please, input the new Device name" },
          {
            min: 4,
            max: 24,
            message:
              "Device name should be at least 4 and not more than 24 characters long",
          },
          {
            pattern: /^[a-zA-Z0-9\-_]*$/,
            message:
              "Device name can contain only alphanumeric characters, hyphens and underscores",
          },
          {
            validator: (rule, val, raise) => {
              let count = (val.match(/[\-_]/g) || []).length;
              if (val.length - count == 0) {
                raise("f#ck");
              }
            },
            message: "Device name can't contain only hyphens and underscores",
          },
        ],
        namespace: [{ required: true, message: "Please select a namespace" }],
        "certificate.pem_data": [
          {
            required: true,
            message: "Please paste or upload device certificate",
          },
        ],
      },
    };
  },
  mounted() {
    this.setDefault();
  },
  methods: {
    setDefault() {
      this.device = {
        name: "",
        tags: [],
        namespace: "",
        enabled: false,
        certificate: { pem_data: "" },
      };
    },
    handleUploadCertificate(file) {
      const reader = new FileReader();

      reader.onload = (e) => {
        this.device.certificate.pem_data = e.target.result;
        this.certificate_tab = "paste";
      };
      reader.readAsText(file);

      return false;
    },
    async handleSubmit() {
      let form = this.$refs["deviceAddForm"];
      let errors = [];
      form.validateField(Object.keys(this.device), (err) => {
        if (err) {
          errors.push(err);
        }
      });
      if (errors.length === 0) this.$emit("add", this.device);
    },
  },
});
</script>

<style>
.ant-form-item-label label {
  font-size: 1rem !important;
}
#certificate-tabs .ant-tabs-tab-active {
  color: var(--primary-color) !important;
  font-weight: 700;
}
#certificate-tabs .ant-tabs-bar {
  border-bottom: none;
}
#certificate-tabs .ant-tabs-nav {
  margin-left: 20% !important;
}
#certificate-tabs .ant-upload.ant-upload-drag p.ant-upload-text {
  color: var(--primary-color);
}
#deviceAddDrawerActionsRow {
  background: var(--primary-color);
}
</style>