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
        <n-form-item label="Title" path="device.title">
          <n-input v-model:value="model.device.title" placeholder="Make it bright" />
        </n-form-item>
        <n-form-item label="Namespace" path="namespace">
          <n-select v-model:value="model.namespace" :options="namespaces" :style="{ minWidth: '15vw' }" />
        </n-form-item>
        <n-form-item label="Enabled" path="device.enabled">
          <n-switch v-model:value="model.device.enabled" />
        </n-form-item>
        <n-form-item label="Tags" path="device.tags">
          <n-dynamic-tags v-model:value="model.device.tags" />
        </n-form-item>
        <n-form-item label="Certificate" path="device.certificate.pem_data">
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
  NButton,
  NDrawer,
  NDrawerContent,
  NIcon,
  NSwitch,
  NSpace,
  NForm,
  NFormItem,
  NInput,
  NDynamicTags,
  NSelect,
  NUpload,
  NUploadDragger,
  NText,
  NAlert,
  useLoadingBar,
} from "naive-ui";
import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { access_lvl_conv } from "@/utils/access";

const AddOutline = defineAsyncComponent(() => import("@vicons/ionicons5/AddOutline"))
const CloudUploadOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CloudUploadOutline"))

const show = ref(false);

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
  return nss.namespaces_list.filter(ns => access_lvl_conv(ns) > 2).map((ns) => ({
    label: `${ns.title} (${shortUUID(ns.uuid)})`,
    value: ns.uuid,
  }));
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
}

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
      console.log(err.response);
      error.value = `${err.response.status}: ${(err.response.data ?? { message: "Unexpected Error" }).message
        }`;
    }
  });
}
</script>