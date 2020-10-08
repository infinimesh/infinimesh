<template>
  <v-container id="footer">
    <v-row justify="center">
      <v-col :sm="9" :md="6" :lg="5" :xl="4" class="footer-text">
        ©2020 —
        <strong>InfiniteDevices GmbH</strong>
        - source code at
        <a href="https://www.github.com/infinimesh/infinimesh" target="_new">
          <strong>GitHub</strong>
        </a>
      </v-col>
      <v-col :sm="3" :md="3" :lg="2" :xl="2">
        <themePicker />
      </v-col>
    </v-row>
  </v-container>
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
  color: var(--footer-text-color) !important;
  text-align: center;
}
#footer strong {
  color: var(--footer-text-color) !important;
}
</style>
