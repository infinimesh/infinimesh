import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { useDevicesStore } from "@/store/devices";

const as = useAppStore();
const nss = useNSStore();
const devs = useDevicesStore();

export const usePluginsStore = defineStore("plugins", {
  state: () => ({
    loading: false,
    plugins: [],

    current: false,

    heights: new Map()
  }),

  actions: {
    async fetchPlugins() {
      this.loading = true;
      const { data } = await as.http.get(
        nss.selected == "all"
          ? "/plugins"
          : `/plugins?namespace=${nss.selected}`
      );

      this.plugins = data.pool;


      this.loading = false;
    },
    async get(uuid) {
      return as.http.get(`/plugins/${uuid}`);
    },
    async create(plugin) {
      return as.http.put("/plugins", plugin);
    },
    async delete(uuid) {
      return as.http.delete("/plugins/" + uuid);
    },
    async update(uuid, data) {
      return as.http.post("/plugins/" + uuid, data);
    },
    height(device) {
      let height = this.heights.get(device);
      if (!height) return '10vh';
      return height + 'px';
    }
  },
});

window.addEventListener('message', ({ origin, data }) => {
  if (!data || !data.type) {
    console.warn("Malformed cross-frame message, skipping. Data:", data, origin);
    return;
  }
  const store = usePluginsStore();
  switch (data.type) {
    case "frame-height":
      console.log(`Setting plugin frame height for ${data.device} to ${data.height}`);
      store.heights.set(data.device, data.height);
      break;
    case "desired":
      console.log(`Received Patch Desired State intent from`, origin);
      console.log(`Device: ${data.device}`, data.state);

      if (!store.current) {
        console.warn(`Plugin ${origin} attempted to patch desired state while not active`);
        return;
      }
      if (!store.current.deviceConf || !store.current.deviceConf.desiredUrl) {
        console.warn(`Current Plugin is either unset or not set for patching desired`);
        return;
      }
      let plugin_origin = (new URL(store.current.deviceConf.desiredUrl)).origin;
      if (plugin_origin != origin) {
        console.warn("Plugin origin is not matching with received message origin. Plugin origin:", plugin_origin, "message origin:", origin);
      }

      devs.patchDesiredState(data.device, data.state, null);

      break;
    default:
      console.warn("Unknown message type", data.type, "from", origin);
  }
});