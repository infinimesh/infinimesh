<template>
  <a-row id="footer" type="flex" justify="center" align="middle">
    <a-col :xs="{ span: 20, offset: 2 }" :sm="18">
      ©2020 —
      <strong>InfiniteDevices GmbH</strong>
      - source code at
      <a href="https://www.github.com/infinimesh/infinimesh" target="_new">
        <strong>GitHub</strong>
      </a>
    </a-col>
    <a-col :xs="1" :sm="2">
      <themePicker />
    </a-col>
  </a-row>
</template>

<script>
import Vue from "vue";
import themePicker from "@/components/generic/themePicker";

const InfinimeshFooter = Vue.component("infinimesh-footer", {
  components: {
    themePicker,
  },
});

export default InfinimeshFooter;
</script>

<style scoped>
#footer {
  width: 90%;
  color: var(--footer-text-color) !important;
  text-align: center;
}
#footer strong {
  color: var(--footer-text-color) !important;
}
</style>
