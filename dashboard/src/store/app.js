import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    token: ""
  }),
  getters: {
    logged_in: (state) => state.token !== ""
  }
})
