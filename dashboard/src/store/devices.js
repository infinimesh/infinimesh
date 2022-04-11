import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { setTransitionHooks } from "vue";

const as = useAppStore();
const nss = useNSStore();

export const useDevicesStore = defineStore("devices", {
  state: () => ({
    loading: false,
    devices: [],
    reported: new Map(),
    desired: new Map(),
    subscribed: [],
  }),

  getters: {
    show_ns: (state) => nss.selected == "all",
    devices_ns_filtered: (state) => {
      let ns = nss.selected;
      let subscribed = new Set(state.subscribed);
      let pool = state.devices
        .map((d) => {
          d.sorter =
            d.enabled + d.accessLevel + d.basicEnabled + subscribed.has(d.uuid);
          return d;
        })
        .sort((a, b) => b.sorter - a.sorter);
      if (ns == "all") {
        return pool;
      }
      return pool.filter((d) => d.namespace == ns);
    },
    device_state: (state) => {
      return (device_id) => {
        return {
          reported: state.reported.get(device_id) ?? {},
          desired: state.desired.get(device_id) ?? {},
        };
      };
    },
    device_subscribed: (state) => {
      return (device_id) => state.subscribed.includes(device_id);
    },
  },

  actions: {
    async fetchDevices() {
      this.loading = true;
      const { data } = await as.http.get("/devices");
      this.devices = data.devices;
      this.loading = false;

      this.getDevicesState(data.devices.map((d) => d.uuid));
    },
    async subscribe(devices) {
      let pool = this.subscribed.concat(devices);

      let token = await this.makeDevicesToken(pool);
      let socket = new WebSocket(`ws://localhost:8000/devices/states/stream`, [
        "Bearer",
        token,
      ]);
      socket.onmessage = (msg) => {
        let response = JSON.parse(msg.data).result;
        if (!response) {
          console.log("Empty response", msg);
          return;
        }

        if (response.reported) {
          this.reported.set(response.device, response.reported);
        }
        if (response.desired) {
          this.desired.set(response.device, response.desired);
        }
      };
      socket.onclose = () => {
        this.subscribed = [];
      };
      socket.onerror = () => {
        this.subscribed = [];
      };
      socket.onopen = () => {
        this.subscribed = pool;
      };
    },
    async makeDevicesToken(pool, post = false) {
      const { data } = await as.http.post("/devices/token", {
        devices: pool,
        post,
      });

      return data.token;
    },
    // pool - array of devices UUIDs
    async getDevicesState(pool, token) {
      if (!token) {
        token = await this.makeDevicesToken(pool);
      }

      const { data } = await as.http.get("/devices/states", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      for (let shadow of data.shadows) {
        this.reported.set(shadow.device, shadow.reported);
        this.desired.set(shadow.device, shadow.desired);
      }
    },
    async patchDesiredState(device, state, bar) {
      bar.start();
      try {
        let token = await this.makeDevicesToken([device], true);
        await as.http.patch(`devices/${device}/state`, state, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.getDevicesState([device], token);
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    async deleteDevice(device, bar) {
      bar.start();
      try {
        await as.http.delete(`devices/${device}`);
        bar.finish();

        this.fetchDevices();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    async createDevice(request, bar) {
      bar.start();
      try {
        await as.http.put(`/devices`, request);

        this.fetchDevices();
        bar.finish();
        return false;
      } catch (e) {
        console.error(e);
        bar.error();
        return e;
      }
    },
    async toggle(uuid, bar) {
      let device;
      for (let dev of this.devices) {
        if (dev.uuid == uuid) {
          device = dev;
          break;
        }
      }
      if (!device) {
        return;
      }

      bar.start();
      device.enabled = null;

      try {
        await as.http.post(`/devices/${uuid}/toggle`);
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
        return;
      }

      this.fetchDevices();
    },
  },
});
