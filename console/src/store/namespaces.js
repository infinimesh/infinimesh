import { useAppStore } from "@/store/app";
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
    }
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh.ns" }],
  },
});
