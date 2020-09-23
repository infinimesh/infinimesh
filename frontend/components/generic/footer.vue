<template>
  <v-container id="footer">
    <v-row justify="center">
      <v-col :sm="9" :md="6" :lg="5" :xl="4">
        ©2020 —
        <strong>InfiniteDevices GmbH</strong>
        - source code at
        <a href="https://www.github.com/infinimesh/infinimesh" target="_new">
          <strong>GitHub</strong>
        </a>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Vue from "vue";

const InfinimeshFooter = Vue.component("infinimesh-footer", {});

export default InfinimeshFooter;
</script>

<style scoped>
#footer {
  color: black;
  text-align: center;
}
</style>
