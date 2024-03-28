<template>
  <n-space vertical v-bind="spaceProps">
    <n-button
      v-for="provider of providers"
      :key="provider"
      :render-icon="oauthIcons(provider)"
      v-bind="buttonProps"
      @click="(type === 'link') ? linkOauth(provider) : oauthLogin(provider)"
    >
      {{ (type === 'link') ? 'Link' : 'Sign in with' }} {{ capitalize(provider) }}
    </n-button>
  </n-space>
</template>

<script setup>
import { ref, defineAsyncComponent, capitalize, h } from 'vue'
import { NButton, NIcon, NSpace, useLoadingBar } from 'naive-ui'
import { useAppStore } from '@/store/app.js'

const githubIcon = defineAsyncComponent(() => import('@vicons/ionicons5/LogoGithub'))
const googleIcon = defineAsyncComponent(() => import('@vicons/ionicons5/LogoGoogle'))

const props = defineProps({
  type: { type: String, default: 'login' },
  buttonProps: {
    type: Object,
    default: () => ({ ghost: true, style: 'width: 100%' })
  },
  spaceProps: {
    type: Object,
    default: () => ({ align: 'center', wrapItem: false })
  }
})
const emits = defineEmits(['success', 'error', 'message'])

const store = useAppStore()
const bar = useLoadingBar()
const providers = ref([])

function oauthIcons(provider) {
  const icons = {
    github: githubIcon,
    google: googleIcon
  }

  return () => h(NIcon, null, { default: () => h(icons[provider]) })
}

async function fetchOauth () {
  try {
    const response = await store.http.get('/oauth/providers')

    providers.value = response.data
  } catch (error) {
    console.error(error)
  }
}

async function oauthLogin (type) {
  emits('success', false)
  emits('error', false)
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

    emits('success', true)
    bar.finish()
  } catch (error) {
    console.error(error)
    emits('error', true)
    bar.error()
  }
}

async function linkOauth (type) {
  bar.start()

  try {
    await store.http.get(
      `/oauth/${type}/login`,
      { params: {
        method: "link",
        state: Math.random().toString(16).slice(2),
        redirect: `https://${location.host}/login`
      } }
    )

    emits('message', { type: 'success', text: 'Done' })
    bar.finish()
  } catch (error) {
    emits('message', {
      type: 'error',
      text: `${error.code}: ${(error.message ?? "Unexpected Error")}`
    })

    console.error(error)
    bar.error()
  }
}

fetchOauth()
</script>