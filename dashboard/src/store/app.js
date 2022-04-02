import { inject } from 'vue'
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    token: "",
    me: {
      title: ""
    },
    namespace_selected: "",
    namespaces: [],
  }),
  getters: {
    logged_in: (state) => state.token !== "",
    http: (state) => {
      return inject('axios').create({
        headers: {
          Authorization: `Bearer ${state.token}`
        }
      })
    }
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
