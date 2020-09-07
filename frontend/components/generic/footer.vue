<template>
  <a-layout-footer id="footer">
    <a-row type="flex" justify="center">
      <a-col :xs="24" :sm="18" :md="12" :lg="10" :xl="8">
        ©2020 —
        <strong>InfniteDevices, GmbH</strong>
        - source code at
        <a
          href="https://www.github.com/infinimesh/infinimesh"
          target="_new"
        >
          <strong>GitHub</strong>
        </a>
      </a-col>
    </a-row>
  </a-layout-footer>
</template>

<script>
import Vue from "vue";

const InfinimeshFooter = Vue.component("infinimesh-footer", {});

export default InfinimeshFooter;
</script>

<style scoped>
#footer {
  color: black;
  margin-top: 15rem;
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  text-align: center;
}
</style>