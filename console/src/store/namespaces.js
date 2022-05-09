import { useAppStore } from "@/store/app";
import { defineStore } from "pinia";

const as = useAppStore();

export const useNSStore = defineStore("namespaces", {
  state: () => ({
    loading: false,
    selected: "",
    namespaces: {},
  }),

  getters: {
    namespaces_list: (state) => {
      return Object.values(state.namespaces)
    }
  },

  actions: {
    async fetchNamespaces() {
      this.loading = true;
      const { data } = await as.http.get("/namespaces");
      this.namespaces = { ...this.namespaces, ...data.namespaces.reduce((r, ns) => { r[ns.uuid] = ns; return r }, {})}
      this.loading = false;
    },
    loadJoins(ns) {
      return as.http.get(`/namespaces/${ns}/joins`)
    },
    join(ns, acc, lvl) {
      return as.http.post(`/namespaces/${ns}/join`, {
        account: acc, access: lvl,
      })
    },
    create(namespace) {
      return as.http.put("/namespaces", namespace);
    }
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh.ns" }],
  },
});
