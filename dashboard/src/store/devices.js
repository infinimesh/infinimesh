import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

const as = useAppStore();
const nss = useNSStore();

export const useDevicesStore = defineStore("devices", {
  state: () => ({
    loading: false,
    devices: [],
    devices_state: new Map(),
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
      return (device_id) => state.devices_state.get(device_id) ?? {};
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

        let curr = this.devices_state.get(response.device);
        let exist = true;
        if (!curr) {
          exist = false;
        }

        if (response.reportedState) {
          curr.reported = response.reportedState;
        }
        if (response.desiredState) {
          curr.desired = response.desiredState;
        }

        if (!exist) {
          this.devices_state.set(response.device, curr);
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

      for (let [uuid, state] of Object.entries(data.pool)) {
        this.devices_state.set(uuid, state);
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
  },
});
