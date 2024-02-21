import { ref, computed } from "vue";
import { defineStore } from "pinia";

import { Struct } from "@bufbuild/protobuf";
import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import {
  DevicesService,
  ShadowService,
} from "infinimesh-proto/build/es/node/node_connect";
import { Level } from "infinimesh-proto/build/es/node/access/access_pb";
import { Shadow } from "infinimesh-proto/build/es/shadow/shadow_pb";
import { Device } from "infinimesh-proto/build/es/node/devices/devices_pb";
import { transport as devicesTransport } from "infinimesh-proto/mocks/es/devices";
import { transport as shadowTransport } from "infinimesh-proto/mocks/es/shadows";

import { useAppStore } from "@/store/app.js";
import { useNSStore } from "@/store/namespaces.js";
import { access_lvl_conv } from "@/utils/access";
import { QueryRequest } from "infinimesh-proto/build/es/node/node_pb";

export const useDevicesStore = defineStore("devices", () => {
  const appStore = useAppStore();
  const namespacesStore = useNSStore();

  const loading = ref(false);
  const devices = ref({});
  const subscribed = ref([]);

  const limit = ref(10);
  const page = ref(1);
  const paginatedDevices = ref([]);
  const paginatedDevicesLoading = ref(false);
  const total = ref(0);

  const reported = ref(new Map());
  const desired = ref(new Map());
  const connection = ref(new Map());

  const devicesApi = computed(() =>
    createPromiseClient(
      DevicesService,
      import.meta.env.VITE_MOCK
        ? devicesTransport
        : createConnectTransport(appStore.transport_options)
    )
  );

  const shadowApi = computed(() =>
    createPromiseClient(
      ShadowService,
      import.meta.env.VITE_MOCK
        ? shadowTransport
        : createConnectTransport(appStore.transport_options)
    )
  );

  const show_ns = computed(() => namespacesStore.selected === "all");
  const devices_ns_filtered = computed(() => {
    const ns = namespacesStore.selected;
    const subscribedSet = new Set(subscribed.value.keys());

    const pool = Object.values(devices.value).map((device) => {
      const { enabled, basicEnabled } = device;
      const level = access_lvl_conv(device);
      const included = subscribedSet.has(device.uuid);

      device.sorter = enabled + level + basicEnabled + included;
      return device;
    });

    pool.sort((a, b) => b.sorter - a.sorter);

    if (ns === "all") return pool;
    return pool.filter((d) => d.access.namespace === ns);
  });

  const device_state = computed(() => (device_id) => ({
    reported: reported.value.get(device_id) ?? {},
    desired: desired.value.get(device_id) ?? {},
    connection: connection.value.get(device_id) ?? {},
  }));
  const device_subscribed = computed(
    () => (device_id) => subscribed.value.includes(device_id)
  );

  async function fetchDevices(state = true, no_cache = false) {
    loading.value = true;
    const data = await devicesApi.value.list();

    if (no_cache) {
      devices.value = data.devices.reduce((result, device) => {
        result[device.uuid] = device;

        return result;
      }, {});
    } else {
      devices.value = {
        ...devices.value,
        ...data.devices.reduce((result, device) => {
          result[device.uuid] = device;

          return result;
        }, {}),
      };
    }

    if (state) getDevicesState(data.devices.map(({ uuid }) => uuid));
    loading.value = false;
  }

  async function fetchDevicesWithPagination(state = true) {
    paginatedDevicesLoading.value = true;
    paginatedDevices.value = [];
    total.value = 0;

    const data = await devicesApi.value.list(
      new QueryRequest({
        offset: (page.value - 1) * limit.value,
        limit: limit.value,
        namespace:
          namespacesStore.selected === "all"
            ? undefined
            : namespacesStore.selected,
      })
    );

    paginatedDevices.value = data.devices;
    total.value = parseInt(data.total);
    paginatedDevicesLoading.value = false;

    state && getDevicesState(data.devices.map((d) => d.uuid));
  }

  async function subscribe(devices) {
    let pool = subscribed.value.concat(devices);

    let token = await makeDevicesToken(pool);
    let socket = new WebSocket(
      `${appStore.base_url.replace("http", "ws")}/devices/states/stream`,
      ["Bearer", token]
    );
    socket.onmessage = (msg) => {
      let response = JSON.parse(msg.data).result;
      if (!response) {
        console.log("Empty response", msg);
        return;
      }

      if (response.reported) {
        if (reported.value.get(response.device)) {
          response.reported.data = {
            ...reported.value.get(response.device).data,
            ...response.reported.data,
          };
        }
        reported.value.set(response.device, response.reported);
      }
      if (response.desired) {
        if (desired.value.get(response.device)) {
          response.desired.data = {
            ...desired.value.get(response.device).data,
            ...response.desired.data,
          };
        }
        desired.value.set(response.device, response.desired);
      }
      if (response.connection) {
        connection.value.set(response.device, response.connection);
      }
    };

    socket.onclose = () => {
      subscribed.value = [];
    };
    socket.onerror = () => {
      subscribed.value = [];
    };
    socket.onopen = () => {
      subscribed.value = pool;
    };
  }

  async function makeDevicesToken(pool, post = false) {
    const level = post ? Level.MGMT : Level.READ;

    const devices = {};
    pool.forEach((uuid) => {
      devices[uuid] = level;
    });

    const data = await devicesApi.value.makeDevicesToken({ devices });
    return data.token;
  }

  // pool - array of devices UUIDs
  async function getDevicesState(pool, token) {
    if (pool.length == 0) return;
    if (!token) {
      token = await makeDevicesToken(pool);
    }

    const data = await shadowApi.value.get(
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    );

    for (const shadow of data.shadows) {
      reported.value.set(shadow.device, shadow.reported);
      desired.value.set(shadow.device, shadow.desired);
      connection.value.set(shadow.device, shadow.connection);
    }
  }

  async function updateDevice(device, patch) {
    if (!patch.title || !patch.tags) {
      throw "Both device Title and Tags must be specified while update";
    }

    try {
      const data = await devicesApi.value.update(
        new Device({ ...patch, uuid: device })
      );

      paginatedDevices.value = paginatedDevices.value.map((dev) => {
        if (dev.uuid === device) {
          return data;
        }
        return dev;
      });
    } catch (error) {
      console.error(error);
      throw `Error Updating Device: ${error.message}`;
    }
  }

  async function updateDeviceConfig(device, config) {
    try {
      const data = await devicesApi.value.update({
        uuid: device,
        config: Struct.fromJson(config),
      });

      paginatedDevices.value = paginatedDevices.value.map((dev) => {
        if (dev.uuid === device) {
          return data;
        }
        return dev;
      });
    } catch (error) {
      console.error(error);
      throw `Error Updating Config: ${error.message}`;
    }
  }

  async function moveDevice(device, namespace) {
    try {
      await devicesApi.value.move({ uuid: device, namespace });

      paginatedDevices.value = paginatedDevices.value.map((dev) => {
        if (dev.uuid === device) {
          dev.access.namespace = namespace;
          return dev;
        }
        return dev;
      });
    } catch (error) {
      console.error(error);
      throw `Error Moving Device: ${error.message}`;
    }
  }

  async function patchDesiredState(device, state, bar) {
    bar?.start();
    try {
      const token = await makeDevicesToken([device], true);
      const data = Struct.fromJson(state);
      const request = new Shadow({ device, desired: { data } });

      await shadowApi.value.patch(request, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      getDevicesState([device], token);
      bar?.finish();
    } catch (e) {
      console.error(e);
      bar?.error();
    }
  }

  async function patchReportedState(device, state, bar) {
    bar.start();
    try {
      const token = await makeDevicesToken([device], true);
      const data = Struct.fromJson(state);
      const request = new Shadow({ device, reported: { data } });

      await shadowApi.value.patch(request, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      getDevicesState([device], token);
      bar.finish();
    } catch (e) {
      console.error(e);
      bar.error();
    }
  }

  async function deleteDevice(device, bar) {
    bar.start();
    try {
      await devicesApi.value.delete({ uuid: device });
      bar.finish();

      fetchDevicesWithPagination();
    } catch (error) {
      console.error(error);
      bar.error();
    }
  }

  async function createDevice(request, bar) {
    bar.start();
    try {
      await devicesApi.value.create(request);

      fetchDevicesWithPagination();
      bar.finish();

      return false;
    } catch (error) {
      console.error(error);
      bar.error();

      return error;
    }
  }

  async function toggle(uuid, bar) {
    const device = paginatedDevices.value.find((dev) => dev.uuid === uuid);
    if (!device) return;

    bar.start();
    device.enabled = null;
    try {
      const data = await devicesApi.value.toggle({ uuid });

      paginatedDevices.value = paginatedDevices.value.map((dev) => {
        if (dev.uuid === uuid) {
          return { ...dev, ...data };
        }
        return dev;
      });
      bar.finish();
    } catch (error) {
      console.error(error);
      bar.error();
    }
  }

  async function toggleBasic(uuid, bar) {
    const device = paginatedDevices.value.find((dev) => dev.uuid === uuid);
    if (!device) return;

    bar.start();
    try {
      const data = await devicesApi.value.toggleBasic({ uuid });

      paginatedDevices.value = paginatedDevices.value.map((dev) => {
        if (dev.uuid === uuid) {
          return { ...dev, ...data };
        }
        return dev;
      });

      bar.finish();
    } catch (error) {
      console.error(error);
      bar.error();
    }
  }

  function fetchJoins(device) {
    return devicesApi.value.joins({ uuid: device });
  }

  async function join(params) {
    return devicesApi.value.join(params);
  }

  return {
    loading,
    devices,
    paginatedDevices,
    limit,
    page,
    total,
    paginatedDevicesLoading,
    subscribed,

    reported,
    desired,
    connection,

    devicesApi,
    shadowApi,
    show_ns,
    devices_ns_filtered,
    device_state,
    device_subscribed,

    fetchDevices,
    fetchDevicesWithPagination,
    fetchJoins,
    subscribe,
    getDevicesState,

    moveDevice,
    updateDevice,
    deleteDevice,
    createDevice,

    makeDevicesToken,
    updateDeviceConfig,
    patchDesiredState,
    patchReportedState,

    join,
    toggle,
    toggleBasic,
  };
});
