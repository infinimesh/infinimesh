<template>
  <n-tooltip :show="min_dpressed" placement="bottom">
    <template #trigger>
      <n-card embedded :bordered="false" hoverable size="huge" title="infinimesh"
        header-style="font-family: 'Exo 2', sans-serif; font-size: 2vh" class="login-card">
        <template #header-extra>
          <n-space>
            <theme-picker />
            <n-button type="info" ghost @click="login">Login</n-button>
          </n-space>
        </template>
        <n-space vertical>
          <n-input v-model:value="username" placeholder="Username" @focus="() => handleDBlock(true)"
            @blur="() => handleDBlock(false)"></n-input>
          <n-input v-model:value="password" type="password" placeholder="Password" @focus="() => handleDBlock(true)"
            @blur="() => handleDBlock(false)" show-password-on="mousedown"></n-input>
          <n-alert :title="error.title" type="error" v-if="error" />

          <n-space justify="center">
            <n-radio-group v-model:value="type" name="credentials_type">
              <n-radio-button value="standard" label="Standard" />
              <n-radio-button value="ldap" label="LDAP" />
            </n-radio-group>
          </n-space>

          <n-alert title="Success! Redirecting..." type="success" v-if="success" />
        </n-space>
      </n-card>
    </template>
    <span v-if="store.dev"> You are now in the developer mode</span>
    <span v-else> You are {{ 10 - dpressed }} d's away from Developer mode</span>
  </n-tooltip>
</template>

<script setup>
import { ref, inject, onMounted, defineAsyncComponent } from "vue";
import {
  NCard,
  NSpace,
  NInput,
  NButton, NRadioButton,
  NAlert, NRadioGroup,
  useLoadingBar, NTooltip
} from "naive-ui";
import { useRoute, useRouter } from "vue-router";
import { useAppStore } from "@/store/app";

const ThemePicker = defineAsyncComponent(() => import("@/components/core/theme-picker.vue"))

const store = useAppStore();
const router = useRouter();

const username = ref("");
const password = ref("");
const type = ref("standard")

const error = ref(false);
const success = ref(false);

const bar = useLoadingBar();

const axios = inject("axios");
async function login() {
  success.value = false;
  error.value = false;
  bar.start();

  const data = {
    auth: {
      type: type.value,
      data: [username.value, password.value],
    },
  }

  if (store.dev) {
    data.inf = true
  }

  axios
    .post(store.base_url + "/token", data)
    .then((res) => {
      success.value = true;
      store.token = res.data.token;
      router.push({ name: "Root" });
      bar.finish();
    })
    .catch((err) => {
      console.error(err);
      bar.error();
      if (err.response.status == 401) {
        error.value = {
          title: "Wrong credentials given",
        };
      }
    });
}

const route = useRoute();

const dpressed = ref(0)
const min_dpressed = ref(false)
const block_dpressed = ref(false)

function handleDBlock(lock) {
  block_dpressed.value = lock
}

onMounted(() => {
  if (route.query.a) {
    try {
      const { token, title, back_token } = JSON.parse(atob(route.query.a))
      store.token = token
      store.me.title = title

      const handle_unload = () => store.token = back_token

      window.addEventListener("beforeunload", handle_unload, false);
      setTimeout(() => {
        handle_unload()
        alert("Token expired. You have been logged out")
        window.close()
      }, 5 * 60 * 1000);

      router.push({ name: "Root" });
    } catch (e) {
      console.error(e)
      error.value = {
        title: "Invalid login token"
      }
    }
  }

  window.addEventListener("keypress", function (e) {
    if (String.fromCharCode(e.keyCode) != 'd' || block_dpressed.value || store.dev) {
      return
    }

    if (block_dpressed.value) {
      return
    }

    dpressed.value += 1
    if (dpressed.value > 2) {
      min_dpressed.value = true

      setTimeout(() => {
        min_dpressed.value = false
        dpressed.value = 0
      }, 3000)
    }
    if (dpressed.value > 9) {
      store.dev = true
    }
  });
})
</script>

<style>
.login-card {
  min-width: 40vw;
}
</style>