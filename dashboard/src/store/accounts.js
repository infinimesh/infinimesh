import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

const as = useAppStore();
const nss = useNSStore();

export const useAccountsStore = defineStore("accounts", {
  state: () => ({
    loading: false,
    accounts: [],
  }),

  getters: {
    accounts_ns_filtered: (state) => {
      let ns = nss.selected;
      let pool = state.accounts
        .map((d) => {
          d.sorter = d.enabled + d.accessLevel;
          return d;
        })
        .sort((a, b) => b.sorter - a.sorter);
      if (ns == "all") {
        return pool;
      }
      return pool.filter((a) => a.namespace == ns);
    },
  },
  actions: {
    async fetchAccounts() {
      this.loading = true;
      const { data } = await as.http.get("/accounts");
      this.accounts = data.accounts;
      this.loading = false;
    },
  },
});
