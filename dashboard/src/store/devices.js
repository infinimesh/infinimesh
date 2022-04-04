import { useAppStore } from "@/store/app";
import { defineStore } from 'pinia'

const as = useAppStore();

export const useDevicesStore = defineStore('devices', {
  state: () => ({
    loading: false,
    devices: [],
  }),

  actions: {
    async fetchDevices() {
      this.loading = true
      const { data } = await as.http.get('/devices');
      this.devices = data.devices;
      this.loading = false
    }
  }
})