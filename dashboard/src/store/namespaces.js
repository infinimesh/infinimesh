import { defineStore } from 'pinia'

export const useNSStore = defineStore('namespaces', {
  state: () => ({
    selected: "",
    namespaces: [],
  }),

  persist: {
    enabled: true,
    strategies: [
      { storage: localStorage, key: 'infinimesh.ns' },
    ],
  },
})