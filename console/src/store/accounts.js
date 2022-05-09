import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

import { access_lvl_conv } from "@/utils/access";

const as = useAppStore();
const nss = useNSStore();

export const useAccountsStore = defineStore("accounts", {
  state: () => ({
    loading: false,
    accounts: {},
  }),

  getters: {
    accounts_ns_filtered: (state) => {
      let ns = nss.selected;
      let pool = Object.values(state.accounts)
        .map((acc) => {
          acc.sorter = acc.enabled + access_lvl_conv(acc);
          return acc;
        })
        .sort((a, b) => b.sorter - a.sorter);
      if (ns == "all") {
        return pool;
      }
      return pool.filter((a) => a.access.namespace == ns);
    },
  },
  actions: {
    async fetchAccounts() {
      this.loading = true;
      const { data } = await as.http.get("/accounts");
      this.accounts = { ...this.accounts, ...data.accounts.reduce((r, ns) => { r[ns.uuid] = ns; return r }, {})};
      this.loading = false;
    },
    async createAccount(request, bar) {
      bar.start();
      try {
        await as.http.put(`/accounts`, request);

        this.fetchAccounts();
        bar.finish();
        return false;
      } catch (e) {
        console.error(e);
        bar.error();
        return e;
      }
    },
    async updateAccount(account, bar) {
      bar.start();
      try {
        await as.http.patch(`/accounts/${account.uuid}`, account);

        this.fetchAccounts();
        bar.finish();
        return false;
      } catch (e) {
        console.error(e);
        bar.error();
        return e;
      }
    },
    async toggle(uuid, bar) {
      bar.start();

      try {
        await as.http.post(`/accounts/${uuid}/toggle`);
        bar.finish();

        this.fetchAccounts();
      } catch (e) {
        console.error(e);
        bar.error();
        return;
      }
    },
    deletables(uuid) {
      return as.http.get(`/accounts/${uuid}/deletables`);
    },
    async deleteAccount(uuid, bar) {
      bar.start();
      try {
        await as.http.delete(`/accounts/${uuid}`);
        bar.finish();

        this.fetchAccounts();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    async setCredentials(uuid, credentials, bar) {
      bar.start();
      try {
        await as.http.post(`/accounts/${uuid}/credentials`, {credentials});
        bar.finish();

        this.fetchAccounts();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    }
  },
});
