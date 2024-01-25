import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { useDevicesStore } from "@/store/devices";
import { createPromiseClient } from "@connectrpc/connect";
import { PluginsService } from "infinimesh-proto/build/es/plugins/plugins_connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const as = useAppStore();
const nss = useNSStore();
const devs = useDevicesStore();

export const usePluginsStore = defineStore("plugins", {
  state: () => ({
    loading: false,
    plugins: [],

    current: false,

    heights: new Map(),
  }),
  getters: {
    plugins_client() {
      return createPromiseClient(
        PluginsService,
        createConnectTransport(as.transport_options)
      );
    },
  },
  actions: {
    async fetchPlugins() {
      this.loading = true;
      const data = await this.plugins_client.list({ namespace: nss.selected });
      this.plugins = data.pool;

      this.loading = false;
    },
    get(uuid) {
      return this.plugins_client.get({ uuid });
    },
    create(plugin) {
      return this.plugins_client.create({
        ...plugin,
        embeddedConf: plugin.embedded_conf,
        deviceConf: plugin.device_conf,
      });
    },
    delete(uuid) {
      return this.plugins_client.delete({ uuid });
    },
    update(uuid, data) {
      return this.plugins_client.update({ ...data, uuid });
    },
    height(device) {
      let height = this.heights.get(device);
      if (!height) return "10vh";
      return height + "px";
    },
  },
});

window.addEventListener("message", ({ origin, data }) => {
  if (!data || !data.type) {
    if (data.source && data.source.includes("vue-devtools-")) return;
    console.warn(
      "Malformed cross-frame message, skipping. Data:",
      data,
      origin
    );
    return;
  }
  const store = usePluginsStore();
  switch (data.type) {
    case "frame-height":
      console.log(
        `Setting plugin frame height for ${data.device} to ${data.height}`
      );
      store.heights.set(data.device, data.height);
      break;
    case "desired":
      console.log(`Received Patch Desired State intent from ${origin}`);
      console.log("Device", data.device, "state", data.state);

      if (!store.current) {
        console.warn(
          `Plugin ${origin} attempted to patch desired state while not active`
        );
        return;
      }
      if (!store.current.deviceConf || !store.current.deviceConf.desiredUrl) {
        console.warn(
          `Current Plugin is either unset or not set for patching desired`
        );
        return;
      }
      let plugin_origin = new URL(store.current.deviceConf.desiredUrl).origin;
      if (plugin_origin != origin) {
        console.warn(
          "Plugin origin is not matching with received message origin. Plugin origin:",
          plugin_origin,
          "message origin:",
          origin
        );
      }

      devs.patchDesiredState(data.device, data.state, null);

      break;
    default:
      console.warn("Unknown message type", data.type, "from", origin);
  }
});
