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
    async sync_me() {
      const { data } = await as.http.get("/accounts/me")
      as.me = { ...as.me, ...data };
    },
    async fetchAccounts(no_cache = false) {
      this.loading = true;

      const { data } = await as.http.get("/accounts");

      if (no_cache) {
        this.accounts = data.accounts.reduce((r, a) => { r[a.uuid] = a; return r; }, {});
      } else {
        this.accounts = { ...this.accounts, ...data.accounts.reduce((r, a) => { r[a.uuid] = a; return r; }, {}) };
      }


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
      if (bar) bar.start();
      try {
        await as.http.patch(`/accounts/${account.uuid}`, account);

        this.fetchAccounts();
        if (bar) bar.finish();
        return false;
      } catch (e) {
        console.error(e);
        if (bar) bar.error();
        return e;
      }
    },
    updateDefaultNamespace(account, ns) {
      let acc = this.accounts[account];
      return this.updateAccount({
        ...acc, defaultNamespace: ns,
      })
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
    async moveAccount(account, namespace) {
      try {
        await as.http.post(`/accounts/${account}/namespace`, { namespace });
        this.accounts[account].access.namespace = namespace;
      } catch (err) {
        console.error(err);
        throw `Error Moving Device: ${err.response.data.message}`;
      }
    },
    async getCredentials(uuid) {
      try {
        const { data } = await as.http.get(`/accounts/${uuid}/credentials`);

        return data;
      } catch (e) {
        console.error(e);
      }
    },
    async setCredentials(uuid, credentials, bar) {
      bar.start();
      try {
        await as.http.post(`/accounts/${uuid}/credentials`, { credentials });
        bar.finish();

        this.fetchAccounts();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    tokenFor(account) {
      return as.http.post(`/token`, {
        uuid: account
      });
    }
  },
});
