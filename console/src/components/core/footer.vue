<template>
  <n-space align="center" justify="space-between" style="min-height: 4vh">
    <span>
      Source code at
      <a href="https://www.github.com/infinimesh/infinimesh" target="_blank">
        <strong>GitHub</strong>
      </a>
      ©2020-2022
    </span>
    <span>
      <span style="font-family: 'Exo 2'">infinimesh</span> - {{ tag }}</span>
    <span></span>
  </n-space>
</template>

<script setup>
import { NSpace } from "naive-ui";

const tag = "development";
</script>