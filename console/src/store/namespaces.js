import { useAppStore } from "@/store/app";
import { join } from "lodash";
import { defineStore } from "pinia";

const as = useAppStore();

export const useNSStore = defineStore("namespaces", {
  state: () => ({
    loading: false,
    selected: "",
    namespaces: [],
  }),

  actions: {
    async fetchNamespaces() {
      this.loading = true;
      const { data } = await as.http.get("/namespaces");
      this.namespaces = data.namespaces;
      this.loading = false;
    },
    loadJoins(ns) {
      return as.http.get(`/namespaces/${ns}/joins`)
    },
    join(ns, acc, lvl) {
      return as.http.post(`/namespaces/${ns}/join`, {
        account: acc, access: lvl,
      })
    }
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh.ns" }],
  },
});
