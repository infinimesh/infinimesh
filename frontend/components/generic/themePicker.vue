<template>
  <a-dropdown placement="topCenter">
    <a class="ant-dropdown-link" @click="(e) => e.preventDefault()">
      <theme-picker-icon />
    </a>
    <a-menu
      slot="overlay"
      :selectedKeys="preference"
      @click="handleThemeChange"
    >
      <a-menu-item :key="mode.value" v-for="mode in modes"
        >{{ mode.title }}
      </a-menu-item>
    </a-menu>
  </a-dropdown>
</template>

<script>
import Vue from "vue";
import themePickerIcon from "./themePickerIcon";

export default {
  components: {
    themePickerIcon,
  },
  computed: {
    preference() {
      return [this.$colorMode.value ? this.$colorMode.value : "unknown"];
    },
  },
  data() {
    return {
      modes: [
        { title: "Default(light)", value: "light" },
        { title: "Night(dark)", value: "dark" },
        { title: "Black and white", value: "black-white" },
      ],
    };
  },
  methods: {
    handleThemeChange(mode) {
      this.$colorMode.preference = mode.key;
    },
  },
};
</script>