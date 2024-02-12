<template>
  <n-space align="center" justify="space-between" style="min-height: 4vh; padding: 0 10px 0 10px">
    <span @click="e => e.preventDefault()">
      <span style="font-family: 'Exo 2'">{{ platform }}</span> - <n-tooltip :show="min_clicked" placement="top">
        <template #trigger>
          <span @click="handler">{{tag}}</span>
        </template>
        <span v-if="store.dev"> You are now in the developer mode</span>
        <span v-else> You are {{ 10 - clicked }} d's away from Developer mode</span>
      </n-tooltip>
    </span>
    <links />
    <span>
      Source code at
      <a href="https://www.github.com/infinimesh/infinimesh" target="_blank">
        <strong>GitHub</strong>
      </a>
      ©2020-2023
    </span>
  </n-space>
</template>

<script setup>
import { h, ref } from "vue"
import { NSpace, NTooltip } from "naive-ui";

import { useAppStore } from "@/store/app";

const store = useAppStore();

const tag = __INFINIMESH_VERSION_TAG__
const platform = PLATFORM_NAME

const clicked = ref(0)
const min_clicked = ref(false)
const block_clicked = ref(false)

function handler() {
  if (block_clicked.value) {
    return
  }

  clicked.value += 1
  if (clicked.value > 2) {
    min_clicked.value = true

    setTimeout(() => {
      min_clicked.value = false
      clicked.value = 0
    }, 3000)
  }
  if (clicked.value > 9) {
    store.dev = true
  }
}

function links() {
  let children = FooterLinks.map(el => {
    return h('span', { class: 'a', onClick: () => window.open(el.href, '_blank')}, el.title)
  })

  return h(NSpace, {}, () => children)
}
</script>

<style>

span.a {
  cursor:pointer;
  font-family: 'Exo 2'
}

span.a:hover {
  text-decoration:underline;
}

</style>