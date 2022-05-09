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
    async fetchNamespaces(no_cache = false) {
      this.loading = true;
      const { data } = await as.http.get("/namespaces");
      if (no_cache) {
        this.namespaces = data.namespaces.reduce((r, ns) => { r[ns.uuid] = ns; return r }, {});
      } else {
        this.namespaces = { ...this.namespaces, ...data.namespaces.reduce((r, ns) => { r[ns.uuid] = ns; return r }, {})}
      }
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
    },
    deletables(uuid) {
      return as.http.get(`/namespaces/${uuid}/deletables`);
    },
    delete(uuid) {
      delete this.namespaces[uuid];
      return as.http.delete(`/namespaces/${uuid}`);
    }
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh.ns" }],
  },
});
