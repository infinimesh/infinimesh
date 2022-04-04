import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces"
import { defineStore } from 'pinia'

const as = useAppStore();
const nss = useNSStore();

export const useDevicesStore = defineStore('devices', {
  state: () => ({
    loading: false,
    devices: [],
    devices_state: new Map()
  }),

  getters: {
    devices_ns_filtered: (state) => {
      let ns = nss.selected;
      if (ns == "all") return state.devices
      return state.devices.filter(d => d.namespace == ns)
    },
    device_state: (state) => {
      return (device_id) => state.devices_state.get(device_id)
    }
  },

  actions: {
    async fetchDevices() {
      this.loading = true
      const { data } = await as.http.get('/devices');
      this.devices = data.devices;
      this.loading = false

      this.getDevicesState(data.devices.map(d => d.uuid))
    },
    // pool - array of devices UUIDs
    async getDevicesState(pool) {
      const { data: res } = await as.http.post('/devices/token', {
        devices: pool, post: false
      })

      let token = res.token;

      const { data } = await as.http.get('/devices/states/all', {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })

      for (let [uuid, state] of Object.entries(data.pool)) {
        this.devices_state.set(uuid, state)
      }
    }
  }
})