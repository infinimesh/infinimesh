import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    token: "",
    me: {
      title: ""
    }
  }),
  getters: {
    logged_in: (state) => state.token !== "",
  },
  actions: {
    logout() {
      console.log(this)
      this.$reset()
    }
  },

  persist: {
    enabled: true,
    strategies: [
      { storage: localStorage, key: 'infinimesh' },
    ],
  },
})
