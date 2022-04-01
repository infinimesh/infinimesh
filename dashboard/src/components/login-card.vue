<template>
  <n-card embedded :bordered="false" hoverable size="huge" title="infinimesh" :header-style="{fontFamily: 'Exo', fontSize: '2vh'}" class="login-card">
    <template #header-extra>
      <n-button type="info" ghost @click="login">Login</n-button>
    </template>
    <n-space vertical>
      <n-input v-model:value="username"
        placeholder="Username"></n-input>
      <n-input
        v-model:value="password"
        type="password"
        placeholder="Password"></n-input>
      <n-alert :title="error.title" type="error" v-if="error" />
      <n-alert title="Success! Redirecting..." type="success" v-if="success" />
    </n-space>
  </n-card>
</template>

<script setup>
import { ref, inject } from "vue";
import { NCard, NSpace, NInput, NButton, NForm, NFormItem, NAlert } from "naive-ui"
import { useRouter } from "vue-router"
import { useAppStore } from "@/store/app";

const store = useAppStore()
const router = useRouter()

const username = ref("");
const password = ref("");

const error = ref(false);
const success = ref(false)

const axios = inject('axios')
async function login() {
  success.value = false
  error.value = false
  
  axios.post('http://localhost:8000/token', {
    auth: {
      type: 'standard',
      data: [username.value, password.value]
    },
  }).then(res => {
    success.value = true
    store.token = res.data.token
    router.push({name: 'Root'})
  }).catch(err => {
    console.log(err)
    if (err.response.status == 401) {
      error.value = {
        title: "Wrong credentials given",
      }
    }
  })
}
</script>

<style>
.login-card {
  min-width: 30vh;
}
</style>