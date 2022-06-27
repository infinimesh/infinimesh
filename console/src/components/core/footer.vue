<template>
  <n-space align="center" justify="space-between" style="min-height: 4vh">
    <span>
      Source code at
      <a href="https://www.github.com/infinimesh/infinimesh" target="_blank">
        <strong>GitHub</strong>
      </a>
      ©2020-2022
    </span>
    <span @click="e => e.preventDefault()">
      <span style="font-family: 'Exo 2'">infinimesh</span> - <n-tooltip :show="min_clicked" placement="top">
        <template #trigger>
          <span @click="handler">{{ tag }}</span>
        </template>
        <span v-if="store.dev"> You are now in the developer mode</span>
        <span v-else> You are {{ 10 - clicked }} d's away from Developer mode</span>
      </n-tooltip></span>
    <span></span>
  </n-space>
</template>

<script setup>
import { ref } from "vue"
import { NSpace, NTooltip } from "naive-ui";

import { useAppStore } from "@/store/app";

const store = useAppStore();

const tag = "development";

const clicked = ref(0)
const min_clicked = ref(false)
const block_clicked = ref(false)

function handler() {
  console.log('clicked')
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
</script>