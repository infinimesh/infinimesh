<template>
  <n-button @click="show = true" type="success" dashed>
    <template #icon>
      <n-icon>
        <add-outline />
      </n-icon>
    </template>
    Create Device
  </n-button>
  <n-drawer v-model:show="show" width="480">
    <n-drawer-content>
      <template #header> Create Device </template>
      <template #footer>
        <n-space justify="end" align="center">
          <n-button type="error" round secondary @click="show = false">Cancel</n-button>
          <n-button type="info" round secondary @click="reset">Reset</n-button>
          <n-button type="warning" round @click="handleSubmit">Submit</n-button>
        </n-space>
      </template>

      <n-form ref="form" :model="model" :rules="rules" label-placement="top">
        <!-- Device Title -->
        <n-form-item label="Title" path="device.title">
          <n-input v-model:value="model.device.title" placeholder="Make it bright" />
        </n-form-item>
        <!-- Device Namespace Selector -->
        <n-form-item label="Namespace" path="namespace">
          <n-select v-model:value="model.namespace" :options="namespaces" :style="{ minWidth: '15vw' }" filterable />
        </n-form-item>
        <!-- Device Enabled -->
        <n-form-item label="Enabled" path="device.enabled" label-placement="left">
          <n-switch v-model:value="model.device.enabled" />
        </n-form-item>
        <!-- Device Tags -->
        <n-form-item label="Tags" path="device.tags">
          <n-dynamic-tags v-model:value="model.device.tags" />
        </n-form-item>

        <!-- Device Credentials -->
        <n-form-item label="Credentials Mode" label-placement="top">
          <n-radio-group v-model:value="mode" name="device_credentials_mode">
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-radio-button value="handsfree" key="handsfree" label="Handsfree" />
              </template>
              Use Authorization code from device to obtain it's pre-installed certificate
            </n-tooltip>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-radio-button value="certificate" key="certificate" :label="'Certificate'" />
              </template>
              Upload your own certificate
            </n-tooltip>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-radio-button value="soft" key="soft" label="Soft(Token)" />
              </template>
              No credentials required, but device will only be able to send data using Device Token.
            </n-tooltip>
          </n-radio-group>
        </n-form-item>
        <!-- Handsfree certificate obtain -->
        <n-form-item label="Code" path="handsfree.code" v-if="mode == 'handsfree'">
          <Code @update:value="(v) => model.handsfree.code = v.code" />
        </n-form-item>
        <!-- Certificate upload -->
        <n-form-item label="Certificate" path="device.certificate.pem_data" v-if="mode == 'certificate'">
          <n-upload v-if="pem_not_uploaded" @before-upload="handleUploadCertificate" accept=".crt,.pem"
            :show-file-list="false">
            <n-upload-dragger>
              <div style="margin-bottom: 12px">
                <n-icon size="48" :depth="3">
                  <cloud-upload-outline />
                </n-icon>
              </div>
              <n-text style="font-size: 16px">
                Click or drag a .crt file to this area to upload
              </n-text>
            </n-upload-dragger>
          </n-upload>
          <n-alert v-else title="Certificate Upload Done" type="success" closable
            @close="model.device.certificate.pem_data = ''">
            Close this alert to upload another certificate
          </n-alert>
        </n-form-item>
        <!-- No Certificate (Soft-device) -->
        <n-alert title="Software Device" type="warning" v-if="mode == 'soft'">
          Credentials are not required, however, the device can only transmit data using a Device Token. This cannot be
          changed later.
        </n-alert>

      </n-form>
      <n-alert title="Error creating Device" type="error" v-if="error">
        {{ error }}
      </n-alert>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup>
import { ref, watch, computed, defineAsyncComponent } from "vue";
import {
  NButton, NDrawer, NDrawerContent, NIcon,
  NSwitch, NSpace, NForm, NFormItem, NTooltip,
  NInput, NDynamicTags, NSelect, NUpload,
  NUploadDragger, NText, NAlert, useLoadingBar,
  NRadioGroup, NRadioButton
} from "naive-ui";
import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { access_lvl_conv } from "@/utils/access";

const AddOutline = defineAsyncComponent(() => import("@vicons/ionicons5/AddOutline"))
const CloudUploadOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloudUploadOutline"))
const Code = defineAsyncComponent(() => import("@/components/devices/register_modal/code.vue"))

const show = ref(false);
const handsfree = ref(false);
const mode = ref('certificate')

watch(
  () => show.value,
  (val) => {
    val && reset();
  }
);

const nss = useNSStore();

function shortUUID(uuid) {
  return uuid.substr(0, 8);
}

const namespaces = computed(() => {
  return nss.namespaces_list.map((ns) => ({
    label: `${ns.title} (${shortUUID(ns.uuid)})`,
    value: ns.uuid,
    disabled: access_lvl_conv(ns) < 3
  })).sort((a, b) => a.disabled - b.disabled);
});

const form = ref();
const model = ref({
  device: {
    title: "",
    enabled: false,
    certificate: {
      pem_data: "",
    },
    tags: [],
  },
  namespace: nss.selected == "all" ? null : nss.selected,
});
const rules = ref({
  device: {
    title: [{ required: true, message: "Please input title" }],
    certificate: {
      pem_data: [{ required: true, message: "Please upload certificate" }],
    },
  },
  namespace: [{ required: true, message: "Please select namespace" }],
  handsfree: {
    code: [{ required: false, message: "Please enter the auth code" }]
  }
});
const store = useDevicesStore();

function reset() {
  model.value = {
    device: {
      title: "",
      enabled: false,
      certificate: {
        pem_data: "",
      },
      tags: [],
    },
    namespace: nss.selected == "all" ? null : nss.selected,
  };
  handsfree.value = false
}

watch(mode, (mode) => {
  model.value.device.certificate = {pem_data: ""}

  if (mode == 'handsfree') {
    rules.value.device.certificate.pem_data[0].required = false
    rules.value.handsfree.code[0].required = true

    model.value.handsfree = { code: "" }

    return
  } else if (mode == 'certificate') {
    rules.value.device.certificate.pem_data[0].required = true
    rules.value.handsfree.code[0].required = false

    delete model.value.handsfree
  } else if (mode == 'soft') {
    rules.value.device.certificate.pem_data[0].required = false
    rules.value.handsfree.code[0].required = false

    delete model.value.handsfree
    delete model.value.device.certificate
  }

  console.log(model.value)
})

const pem_not_uploaded = computed(
  () => model.value.device.certificate.pem_data == ""
);

function handleUploadCertificate({ file }) {
  const reader = new FileReader();

  reader.onload = (e) => {
    model.value.device.certificate.pem_data = e.target.result;
  };
  reader.readAsText(file.file);

  return false;
}

const error = ref(false);
const bar = useLoadingBar();
function handleSubmit() {
  error.value = false;
  form.value.validate(async (errors) => {
    if (errors) {
      return;
    }
    let err = await store.createDevice(model.value, bar);
    if (!err) {
      show.value = false;
    } else {
      console.log(err);
      error.value = `${err.code}: ${err.message ?? "Unexpected Error"}`;
    }
  });
}
</script>