<template>
  <a-drawer
    title="Create New Device"
    :visible="active"
    width="30%"
    :headerStyle="{ fontSize: '4rem' }"
    @close="$emit('cancel')"
  >
    <a-form-model
      :model="device"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-model-item prop="name" label="Name">
        <a-input v-model="device.name" />
      </a-form-model-item>
      <a-form-model-item label="Namespace" prop="namespace">
        <a-select
          v-model="device.namespace"
          placeholder="please select device namespace"
        >
          <a-select-option :value="ns.id" :key="ns.id" v-for="ns in namespaces">
            {{ ns.name }}
          </a-select-option>
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
    </a-form-model>
  </a-drawer>
</template>

<script>
import Vue from "vue";

export default Vue.component("device-add", {
  props: {
    active: {
      required: true
    }
  },
  computed: {
    namespaces: {
      get() {
        return this.$store.state.devices.namespaces;
      }
    }
  },
  watch: {
    active: "setDefault"
  },
  data() {
    return {
      device: {},
      rules: {
        name: [
          { required: true, message: "Please input the new Device name" },
          {
            min: 4,
            max: 24,
            message:
              "Device name should be at least 4 and not more than 24 characters long"
          },
          {
            pattern: /^[a-zA-Z0-9\-_]*$/,
            message:
              "Device name can contain only alphanumeric characters, hyphens and underscores"
          },
          {
            validator: (rule, val, raise) => {
              console.log(val);
              let count = (val.match(/[\-_]/g) || []).length;
              if (val.length - count == 0) {
                raise("f#ck");
              }
            },
            message: "Device name can't contain only hyphens and underscores"
          }
        ],
        namespace: [{ required: true, message: "Please select a namespace" }]
      }
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
        certificate: ""
      };
    }
  }
});
</script>

<style>
.ant-form-item-label label {
  font-size: 1rem !important;
}
</style>
