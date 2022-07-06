import { inject } from "vue";
import { defineStore } from "pinia";

import { check_token_expired, check_offline } from "@/utils/access";

export const baseURL = import.meta.env.DEV
  ? "http://api.infinimesh.local"
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

      const store = this
      function err_check(err) {
        check_token_expired(err, store)
      }

      instance.interceptors.response.use((r) => r, err_check)
      return instance
    },
  },
  actions: {
    logout() {
      this.$reset();
    },
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh" }],
  },
});
