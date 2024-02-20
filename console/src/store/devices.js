import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import {
  DevicesService,
  ShadowService,
} from "infinimesh-proto/build/es/node/node_connect";

import { access_lvl_conv } from "@/utils/access";

import { Level } from "infinimesh-proto/build/es/node/access/access_pb";
import { Shadow } from "infinimesh-proto/build/es/shadow/shadow_pb";
import {
  DevicesTokenRequest,
  QueryRequest,
} from "infinimesh-proto/build/es/node/node_pb";
import { Device } from "infinimesh-proto/build/es/node/devices/devices_pb";
import { Struct } from "@bufbuild/protobuf";

const as = useAppStore();
const nss = useNSStore();

export const useDevicesStore = defineStore("devices", {
  state: () => ({
    loading: false,
    devices: {},
    subscribed: [],

    limit: 10,
    page: 1,
    paginatedDevices: [],
    paginatedDevicesLoading: false,
    total: 0,

    reported: new Map(),
    desired: new Map(),
    connection: new Map(),
  }),

  getters: {
    devices_client() {
      return createPromiseClient(
        DevicesService,
        createConnectTransport(as.transport_options)
      );
    },
    shadow_client() {
      return createPromiseClient(
        ShadowService,
        createConnectTransport({ ...as.transport_options, interceptors: [] })
      );
    },
    show_ns: (state) => nss.selected == "all",
    devices_ns_filtered: (state) => {
      let ns = nss.selected;
      let subscribed = new Set(state.subscribed);
      let pool = Object.values(state.devices)
        .map((d) => {
          d.sorter =
            d.enabled +
            access_lvl_conv(d) +
            d.basicEnabled +
            subscribed.has(d.uuid);
          return d;
        })
        .sort((a, b) => b.sorter - a.sorter);
      if (ns == "all") {
        return pool;
      }
      return pool.filter((d) => d.access.namespace == ns);
    },
    device_state: (state) => {
      return (device_id) => {
        return {
          reported: state.reported.get(device_id) ?? {},
          desired: state.desired.get(device_id) ?? {},
          connection: state.connection.get(device_id) ?? {},
        };
      };
    },
    device_subscribed: (state) => {
      return (device_id) => state.subscribed.includes(device_id);
    },
  },

  actions: {
    async fetchDevices(state = true, no_cache = false) {
      this.loading = true;

      const data = await this.devices_client.list();

      if (no_cache) {
        this.devices = data.devices.reduce((r, d) => {
          r[d.uuid] = d;
          return r;
        }, {});
      } else {
        this.devices = {
          ...this.devices,
          ...data.devices.reduce((r, d) => {
            r[d.uuid] = d;
            return r;
          }, {}),
        };
      }

      if (state) this.getDevicesState(data.devices.map((d) => d.uuid));
      this.loading = false;
    },
    async fetchDevicesWithPagination(state = true) {
      this.paginatedDevicesLoading = true;
      this.paginatedDevices = [];
      this.total = 0;

      const data = await this.devices_client.list(
        new QueryRequest({
          offset: (this.page - 1) * this.limit,
          limit: this.limit,
          namespace: nss.selected === "all" ? undefined : nss.selected,
        })
      );

      this.paginatedDevices = data.devices;
      this.total = parseInt(data.total);
      this.paginatedDevicesLoading = false;

      state && this.getDevicesState(data.devices.map((d) => d.uuid));
    },
    async subscribe(devices) {
      let pool = this.subscribed.concat(devices);

      let token = await this.makeDevicesToken(pool);
      let socket = new WebSocket(
        `${as.base_url.replace("http", "ws")}/devices/states/stream`,
        ["Bearer", token]
      );
      socket.onmessage = (msg) => {
        let response = JSON.parse(msg.data).result;
        if (!response) {
          console.log("Empty response", msg);
          return;
        }

        if (response.reported) {
          if (this.reported.get(response.device)) {
            response.reported.data = {
              ...this.reported.get(response.device).data,
              ...response.reported.data,
            };
          }
          this.reported.set(response.device, response.reported);
        }
        if (response.desired) {
          if (this.desired.get(response.device)) {
            response.desired.data = {
              ...this.desired.get(response.device).data,
              ...response.desired.data,
            };
          }
          this.desired.set(response.device, response.desired);
        }
        if (response.connection) {
          this.connection.set(response.device, response.connection);
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
      const level = post ? Level.MGMT : Level.READ;

      const devices = {};
      pool.forEach((uuid) => {
        devices[uuid] = level;
      });

      const data = await this.devices_client.makeDevicesToken({ devices });

      return data.token;
    },
    // pool - array of devices UUIDs
    async getDevicesState(pool, token) {
      if (pool.length == 0) return;
      if (!token) {
        token = await this.makeDevicesToken(pool);
      }

      const data = await this.shadow_client.get(
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      for (let shadow of data.shadows) {
        this.reported.set(shadow.device, shadow.reported);
        this.desired.set(shadow.device, shadow.desired);
        this.connection.set(shadow.device, shadow.connection);
      }
    },
    async updateDevice(device, patch) {
      if (!patch.title || !patch.tags)
        throw "Both device Title and Tags must be specified while update";
      try {
        const data = await this.devices_client.update(
          new Device({
            ...patch,
            uuid: device,
            config: undefined,
          })
        );
        this.paginatedDevices = this.paginatedDevices.map((d) => {
          if (d.uuid === device) {
            return data;
          }
          return d;
        });
      } catch (err) {
        console.error(err);
        throw `Error Updating Device: ${err.message}`;
      }
    },
    async updateDeviceConfig(device, config) {
      try {
        const data = await this.devices_client.patchConfig({
          uuid: device,
          config: new Struct().fromJson(config),
        });

        this.paginatedDevices = this.paginatedDevices.map((d) => {
          if (d.uuid === device) {
            return data;
          }
          return d;
        });
      } catch (err) {
        console.error(err);
        throw `Error Updating Config: ${err.message}`;
      }
    },
    async moveDevice(device, namespace) {
      try {
        await this.devices_client.move({ uuid: device, namespace });

        this.paginatedDevices = this.paginatedDevices.map((d) => {
          if (d.uuid === device) {
            d.access.namespace = namespace;
          }
          return d;
        });
      } catch (err) {
        console.error(err);
        throw `Error Moving Device: ${err.message}`;
      }
    },
    async patchDesiredState(device, state, bar) {
      if (bar) bar.start();
      try {
        let token = await this.makeDevicesToken([device], true);
        const data = Struct.fromJson(state);
        const request = new Shadow({
          device,
          desired: { data },
        });
        await this.shadow_client.patch(request, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.getDevicesState([device], token);
        if (bar) bar.finish();
      } catch (e) {
        console.error(e);
        if (bar) bar.error();
      }
    },
    async patchReportedState(device, state, bar) {
      bar.start();
      try {
        let token = await this.makeDevicesToken([device], true);
        const data = Struct.fromJson(state);
        const request = new Shadow({
          device,
          reported: { data },
        });
        await this.shadow_client.patch(request, {
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
        await this.devices_client.delete({ uuid: device });
        bar.finish();

        this.fetchDevicesWithPagination(false);
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    async createDevice(request, bar) {
      bar.start();
      try {
        await this.devices_client.create(request);

        this.fetchDevicesWithPagination();
        bar.finish();
        return false;
      } catch (e) {
        console.error(e);
        bar.error();
        return e;
      }
    },
    async toggle(uuid, bar) {
      let device = this.paginatedDevices.find((d) => d.uuid === uuid);
      if (!device) {
        return;
      }

      bar.start();
      device.enabled = null;

      try {
        const data = await this.devices_client.toggle({ uuid });

        this.paginatedDevices = this.paginatedDevices.map((d) => {
          if (d.uuid === uuid) {
            return { ...d, ...data };
          }
          return d;
        });
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
        return;
      }
    },
    async toggle_basic(uuid, bar) {
      let device = this.paginatedDevices.find((d) => d.uuid === uuid);
      if (!device) {
        return;
      }

      bar.start();

      try {
        const data = await this.devices_client.toggleBasic({ uuid });

        this.paginatedDevices = this.paginatedDevices.map((d) => {
          if (d.uuid === uuid) {
            return { ...d, ...data };
          }
          return d;
        });

        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
        return;
      }
    },
    async fetchJoins(device) {
      const data = await this.devices_client.joins({ uuid: device });
      return data;
    },
    async join(params) {
      return this.devices_client.join(params);
    },
  },
});
