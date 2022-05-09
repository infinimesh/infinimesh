<template>
  <n-button @click="show = true" type="success" dashed>
    <template #icon>
      <n-icon>
        <add-outline />
      </n-icon>
    </template>
    Create Account
  </n-button>
  <n-drawer v-model:show="show" width="480">
    <n-drawer-content>
      <template #header> Create Account </template>
      <template #footer>
        <n-space justify="end" align="center">
          <n-button type="error" round secondary @click="show = false">Cancel</n-button>
          <n-button type="info" round secondary @click="reset">Reset</n-button>
          <n-button type="warning" round @click="handleSubmit">Submit</n-button>
        </n-space>
      </template>
      <n-form ref="form" :model="model" label-placement="top">
        <n-form-item label="Title" path="account.title">
          <n-input v-model:value="model.account.title" placeholder="How should we call you?" />
        </n-form-item>
        <n-form-item label="Enabled" path="account.enabled">
          <n-switch v-model:value="model.account.enabled" />
        </n-form-item>
        <n-form-item label="Namespace" path="namespace">
          <n-select v-model:value="model.namespace" :options="namespaces" :style="{ minWidth: '15vw' }" />
        </n-form-item>
        Credentials:
        <n-tabs v-model="model.credentials.type">
          <n-tab-pane name="standard" display-directive="if" tab="Standard(user/pass)">
            <n-form-item label="Username" path="credentials.data[0]">
              <n-input v-model:value="model.credentials.data[0]" />
            </n-form-item>
            <n-form-item label="Password" path="credentials.data[1]">
              <n-input v-model:value="model.credentials.data[1]" type="password" show-password-on="click">
                <template #password-visible-icon>
                  <n-icon :size="16" :component="EyeOffOutline" />
                </template>
                <template #password-invisible-icon>
                  <n-icon :size="16" :component="EyeOutline" />
                </template>
              </n-input>
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup>
import { ref, watch, computed } from "vue";
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
  NSelect,
  NTabs,
  NTabPane,
  useLoadingBar,
} from "naive-ui";
import { AddOutline, EyeOffOutline, EyeOutline } from "@vicons/ionicons5";
import { useNSStore } from "@/store/namespaces";
import { useAccountsStore } from "@/store/accounts";

const show = ref(false);
watch(
  () => show.value,
  (val) => {
    val && reset();
  }
);

const nss = useNSStore();
const namespaces = computed(() => {
  return nss.namespaces_list.map((ns) => ({
    label: ns.title,
    value: ns.uuid,
  }));
});

const form = ref();
const model = ref({
  account: {
    title: "",
    enabled: false,
  },
  namespace: nss.selected == "all" ? null : nss.selected,
  credentials: {
    type: "standard",
    data: ["", ""],
  }
});

watch(() => model.value.account.title, (val) => {
  if (model.value.credentials.type === "standard") {
    model.value.credentials.data[0] = val
  }
})

function reset() {
  model.value = {
    account: {
      title: "",
      enabled: false,
    },
    namespace: nss.selected == "all" ? null : nss.selected,
    credentials: {
      type: "standard",
      data: ["", ""],
    }
  }
}

const store = useAccountsStore()
const error = ref(false);
const bar = useLoadingBar();
function handleSubmit() {
  error.value = false;
  form.value.validate(async (errors) => {
    if (errors) {
      return;
    }
    let err = await store.createAccount(model.value, bar);
    if (!err) {
      show.value = false;
    } else {
      console.log(err.response);
      error.value = `${err.response.status}: ${(err.response.data ?? { message: "Unexpected Error" }).message}`;
    }
  })
}
</script>