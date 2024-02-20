import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { useDevicesStore } from "@/store/devices";
import { createPromiseClient } from "@connectrpc/connect";
import { PluginsService } from "infinimesh-proto/build/es/plugins/plugins_connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ref, computed } from "vue";

import { transport } from "infinimesh-proto/mocks/es/plugins";


export const usePluginsStore = defineStore("plugins", () => {
  const as = useAppStore();
  const nss = useNSStore();

  const loading = ref(false);
  const current = ref(false);
  const plugins = ref([]);
  const heights = ref(new Map());

  const plugins_client = computed(() =>
    createPromiseClient(
      PluginsService,
      import.meta.env.VITE_MOCK
        ? transport
        : createConnectTransport(as.transport_options)
    )
  );

  async function fetchPlugins() {
    loading.value = true;
    const data = await plugins_client.value.list({ namespace: nss.selected });
    plugins.value = data.pool;

    loading.value = false;
  }

  function get(uuid) {
    return plugins_client.value.get({ uuid });
  }

  function create(plugin) {
    return plugins_client.value.create({
      ...plugin,
      embeddedConf: plugin.embedded_conf,
      deviceConf: plugin.device_conf,
    });
  }

  function deletePlugin(uuid) {
    return plugins_client.value.delete({ uuid });
  }

  function update(uuid, data) {
    return plugins_client.value.update({ ...data, uuid });
  }

  function height(device) {
    let height = heights.value.get(device);
    if (!height) return "10vh";
    return height + "px";
  }

  return {
    loading,
    current,
    fetchPlugins,
    get,
    heights,
    plugins_client,
    delete: deletePlugin,
    create,
    update,
    height,
  };
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

      const devs = useDevicesStore();
      devs.patchDesiredState(data.device, data.state, null);

      break;
    default:
      console.warn("Unknown message type", data.type, "from", origin);
  }
});
