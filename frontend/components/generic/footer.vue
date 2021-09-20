<template>
  <a-row id="footer" type="flex" justify="center" align="middle">
    <a-col :xs="{ span: 20, offset: 2 }" :sm="18">
      <a v-if="release" :href="release.html_url" target="_blank"
        ><strong> infinimesh {{ release.tag_name }}</strong></a
      >
      ©2021 —
      <strong>The infinimesh team</strong>
      - source code at
      <a href="https://www.github.com/infinimesh/infinimesh" target="_blank">
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
  computed: {
    release: {
      deep: true,
      get() {
        return this.$store.getters["window/release"];
      },
    },
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
