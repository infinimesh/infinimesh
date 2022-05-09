<template>
  <n-card embedded :bordered="false" hoverable size="huge" title="infinimesh"
    header-style="font-family: 'Exo 2', sans-serif; font-size: 2vh" class="login-card">
    <template #header-extra>
      <n-space>
        <theme-picker />
        <n-button type="info" ghost @click="login">Login</n-button>
      </n-space>
    </template>
    <n-space vertical>
      <n-input v-model:value="username" placeholder="Username"></n-input>
      <n-input v-model:value="password" type="password" placeholder="Password"></n-input>
      <n-alert :title="error.title" type="error" v-if="error" />
      <n-alert title="Success! Redirecting..." type="success" v-if="success" />
    </n-space>
  </n-card>
</template>

<script setup>
import { ref, inject, onMounted } from "vue";
import {
  NCard,
  NSpace,
  NInput,
  NButton,
  NAlert,
  useLoadingBar,
} from "naive-ui";
import { useRoute, useRouter } from "vue-router";
import { useAppStore } from "@/store/app";
import ThemePicker from "@/components/core/theme-picker.vue";

const store = useAppStore();
const router = useRouter();

const username = ref("");
const password = ref("");

const error = ref(false);
const success = ref(false);

const bar = useLoadingBar();

const axios = inject("axios");
async function login() {
  success.value = false;
  error.value = false;
  bar.start();

  axios
    .post(store.base_url + "/token", {
      auth: {
        type: "standard",
        data: [username.value, password.value],
      },
    })
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
})
</script>

<style>
.login-card {
  min-width: 30vh;
}
</style>