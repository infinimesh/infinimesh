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

  actions: {
    async fetchAccounts() {
      this.loading = true;
      const { data } = await as.http.get("/accounts");
      this.accounts = data.accounts;
      this.loading = false;
    },
  },
});
