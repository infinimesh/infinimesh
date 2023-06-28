import { inject, nextTick } from "vue";
import { defineStore } from "pinia";

import { check_token_expired, check_offline } from "@/utils/access";

export const baseURL =
  import.meta.env.DEV ? "http://api.infinimesh.local" // jshint ignore:line
    : window.location.origin.replace("console.", "api.");

export const useAppStore = defineStore("app", {
  state: () => ({
    token: "",
    me: {
      title: "",
    },
    theme: "dark",
    theme_pick: "system",
    console_services: {},
    dev: false,

    current_thing: false,
    notify: null
  }),
  getters: {
    base_url: () => baseURL,
    logged_in: (state) => state.token !== "",
    http(state) {
      const instance = inject("axios").create({
        baseURL,
        headers: {
          Authorization: `Bearer ${state.token}`,
        },
      });

      const store = this;
      function err_check(err) {
        check_token_expired(err, store);
        check_offline(err, store);
        return Promise.reject(err);
      }

      instance.interceptors.response.use((r) => r, err_check);
      return instance;
    },
  },
  actions: {
    logout(msg = false) {
      this.$reset();

      let query = {};
      if (msg) {
        query.msg = btoa(JSON.stringify(msg));
      }

      this.$router.push({ name: 'Login', query });
    },
    offline() {
      this.$router.push({ name: 'Offline' });
    }
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh" }],
  },
});