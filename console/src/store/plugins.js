import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

import { access_lvl_conv, check_token_expired } from "@/utils/access";

const as = useAppStore();
const nss = useNSStore();

export const usePluginsStore = defineStore("plugins", {
  state: () => ({
    loading: false,
    plugins: [],
  }),

  actions: {
    async fetchPlugins() {
      this.loading = true;
      try {
        const { data } = await as.http.get(
          nss.selected == "all"
            ? "/plugins"
            : `/plugins?namespace=${nss.selected}`
        );

        this.plugins = data.pool;
      } catch (e) {
        check_token_expired(e, as);
      }

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
  },
});
