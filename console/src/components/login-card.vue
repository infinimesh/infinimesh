<template>
  <n-tooltip :show="min_dpressed" placement="bottom">
    <template #trigger>
      <n-card embedded :bordered="false" hoverable size="huge" :title="platform"
        header-style="font-family: 'Exo 2', sans-serif; font-size: 2vh" class="login-card">
        <template #header-extra>
          <n-space>
            <theme-picker />
            <n-button v-if="type !== 'oauth'" type="info" ghost @click="login">Login</n-button>
          </n-space>
        </template>
        <n-space vertical>
          <template v-if="type !== 'oauth'">
            <n-input
              v-model:value="username"
              placeholder="Username"
              @focus="handleDBlock(true)"
              @blur="handleDBlock(false)"
            />
            <n-input
              v-model:value="password"
              type="password"
              placeholder="Password"
              show-password-on="mousedown"
              @focus="handleDBlock(true)"
              @blur="handleDBlock(false)"
            />
            <n-alert :title="error.title" type="error" v-if="error" />
          </template>
          <n-button v-else ghost type="primary" @click="oauthLogin('github')">GitHub</n-button>

          <n-space justify="center">
            <n-radio-group v-model:value="type" name="credentials_type">
              <n-radio-button value="standard" label="Standard" />
              <n-radio-button value="ldap" label="LDAP" />
              <n-radio-button value="oauth" label="OAuth 2.0" />
            </n-radio-group>
          </n-space>

          <n-alert title="Success! Redirecting..." type="success" v-if="success" />

          <n-alert v-if="alert" :title="alert.title" :type="alert.type" closable>
          {{  alert.description }}  
          </n-alert>
        </n-space>
      </n-card>
    </template>
    <span v-if="store.dev"> You are now in the developer mode</span>
    <span v-else> You are {{ 10 - dpressed }} d's away from Developer mode</span>
  </n-tooltip>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent } from "vue";
import {
  NCard, NSpace, NInput,
  NButton, NRadioButton,
  NAlert, NRadioGroup,
  useLoadingBar, NTooltip
} from "naive-ui";
import { useRoute, useRouter } from "vue-router";
import { useAppStore } from "@/store/app";
import { useAccountsStore } from "@/store/accounts";

import { UAParser } from "ua-parser-js"

const ThemePicker = defineAsyncComponent(() => import("@/components/core/theme-picker.vue"))

const platform = PLATFORM_NAME

const store = useAppStore();
const accStore = useAccountsStore()
const router = useRouter();

const username = ref("");
const password = ref("");
const type = ref("standard")

const error = ref(false);
const success = ref(false);
const alert = ref(false);

const bar = useLoadingBar();

async function oauthLogin (type) {
  success.value = false
  error.value = false
  alert.value = false
  bar.start()

  try {
    await store.http.get(
      `/oauth/${type}/login`,
      { params: {
        method: 'sign_in',
        state: Math.random().toString(16).slice(2),
        redirect: `https://${location.host}/login`
      } }
    )

    success.value = true
    bar.finish()
  } catch (err) {
    console.error(err)
    error.value = true
    bar.error()
  }
}

onMounted(() => {
  if (route.query.token) {
    store.token = route.query.token
    router.replace({ name: 'root' }).then(() =>
      location.reload()
    )
  }
})

async function login() {
  success.value = false;
  error.value = false;
  alert.value = false;
  bar.start();

  const data = {
    auth: {
      type: type.value,
      data: [username.value, password.value],
    },
  };

  if (store.dev) data.inf = true;

  let res = {};
  try {
    res = new UAParser(navigator.userAgent).getResult()
  } catch (e) {
    console.warn("Failed to get user agent", e)
  }
  data.client = `Console | ${res.os?.name ?? 'Unknown'} | ${res.browser?.name ?? 'Unknown'}`;

  try {
    const { token } = await accStore.token(data)

    success.value = true;
    store.token = token;
    router.push({ name: "Root" });
    bar.finish();
  } catch (e) {
    console.error(e);
    bar.error();

    if (e.message.includes("Wrong credentials given")) {
      error.value = { title: "Wrong credentials given" };
    }
  }
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

  if (route.query.msg) {
    alert.value = JSON.parse(atob(route.query.msg))
  }
})
</script>

<style>
.login-card {
  min-width: 30vw;
  max-width: 90vw;
}
</style>