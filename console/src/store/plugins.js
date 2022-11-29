import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

const as = useAppStore();
const nss = useNSStore();

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

window.addEventListener('message', ({ data }) => {
  if (data.type == "frame-height") {
    console.log(`Setting plugin frame height for ${data.device} to ${data.height}`)
    usePluginsStore().heights.set(data.device, data.height)
  }
});