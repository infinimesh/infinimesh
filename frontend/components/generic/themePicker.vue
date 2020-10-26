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
  mounted() {
    let d = new Date();
    if (
      (d.getDate() > 27 && d.getMonth() == 9) ||
      (d.getDate() < 2 && d.getMonth() == 10)
    ) {
      var style = document.createElement("style");
      style.type = "text/css";
      style.innerHTML = atob(
        "LmhhbGxvd2Vlbi1tb2RlIGJvZHkgewogIC0tYmFja2dyb3VuZC1jb2xvcjogI2ZmNmExZjsKICAtLWZvb3Rlci1iYWNrZ3JvdW5kOiAjZmY4NjFmOwogIC0tZm9vdGVyLXRleHQtY29sb3I6IGJsYWNrOwogIC0tcHJpbWFyeS1jb2xvcjogYmxhY2s7Cn0="
      );
      document.getElementsByTagName("head")[0].appendChild(style);
      this.modes.push({
        title: "🦇",
        value: atob("aGFsbG93ZWVu"),
      });
      this.$colorMode.preference = atob("aGFsbG93ZWVu");
    }
  },
  methods: {
    handleThemeChange(mode) {
      this.$colorMode.preference = mode.key;
    },
  },
};
</script>