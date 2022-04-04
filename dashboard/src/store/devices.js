import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces"
import { defineStore } from 'pinia'

const as = useAppStore();
const nss = useNSStore();

export const useDevicesStore = defineStore('devices', {
  state: () => ({
    loading: false,
    devices: [],
  }),

  getters: {
    devices_ns_filtered: (state) => {
      let ns = nss.selected;
      if (ns == "all") return state.devices
      return state.devices.filter(d => d.namespace == ns)
    }
  },

  actions: {
    async fetchDevices() {
      this.loading = true
      const { data } = await as.http.get('/devices');
      this.devices = data.devices;
      this.loading = false
    }
  }
})