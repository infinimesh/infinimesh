import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";

import { access_lvl_conv } from "@/utils/access";
import { createPromiseClient } from "@connectrpc/connect";
import { AccountsService } from "infinimesh-proto/build/es/node/node_connect"
import { Account, CreateRequest } from "infinimesh-proto/build/es/node/accounts/accounts_pb";
import { EmptyMessage, MoveRequest, SetCredentialsRequest, TokenRequest } from "infinimesh-proto/build/es/node/node_pb";

const as = useAppStore();
const nss = useNSStore();

export const useAccountsStore = defineStore("accounts", {
  state: () => ({
    loading: false,
    accounts: {},
    accountsApi: createPromiseClient(AccountsService, as.transport)
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
      const data = await this.accountsApi.get({ uuid: "me" })

      as.me = { ...as.me, ...data };
    },
    async fetchAccounts(no_cache = false) {
      this.loading = true;

      const { accounts } = await this.accountsApi.list(new EmptyMessage());

      if (no_cache) {
        this.accounts = accounts.reduce((result, account) => {
          result[account.uuid] = account;

          return result;
        }, {});
      } else {
        this.accounts = {
          ...this.accounts,
          ...accounts.reduce((result, account) => {
            result[account.uuid] = account;

            return result;
          }, {})
        };
      }

      this.loading = false;
    },
    async createAccount(request, bar) {
      bar.start();
      try {
        await this.accountsApi.create(new CreateRequest(request));

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
        if (!account.config) account.config = {}
        const result = new Account(account);

        result.config = result.config.fromJson(account.config)
        await this.accountsApi.update(result);

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
        await this.accountsApi.toggle(new Account(this.accounts[uuid]));

        this.fetchAccounts();
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
        return;
      }
    },
    deletables(uuid) {
      return this.accountsApi.deletables(new Account(this.accounts[uuid]));
    },
    async deleteAccount(uuid, bar) {
      bar.start();
      try {
        await this.accountsApi.delete(new Account(this.accounts[uuid]));

        this.fetchAccounts();
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    async moveAccount(uuid, namespace) {
      try {
        await this.accountsApi.move(new MoveRequest({ uuid, namespace }));

        this.accounts[uuid].access.namespace = namespace;
      } catch (err) {
        console.error(err);
        throw `Error Moving Device: ${err.message}`;
      }
    },
    async getCredentials(uuid) {
      try {
        return await this.accountsApi.getCredentials({ uuid });
      } catch (e) {
        console.error(e);
      }
    },
    async setCredentials(uuid, credentials, bar) {
      bar.start();
      try {
        await this.accountsApi.setCredentials(
          new SetCredentialsRequest({ uuid, credentials })
        );

        this.fetchAccounts();
        bar.finish();
      } catch (e) {
        console.error(e);
        bar.error();
      }
    },
    token (data) {
      try {
        return this.accountsApi.token(new TokenRequest(data))
      } catch (e) {
        console.error(e); 
      }
    },
    tokenFor(account, exp = 0) {
      let res = {};
      try {
        res = new UAParser(navigator.userAgent).getResult()
      } catch (e) {
        console.warn("Failed to get user agent", e)
      }

      return this.token({
        uuid: account,
        client: `Console Admin | ${res.os?.name ?? 'Unknown'} | ${res.browser?.name ?? 'Unknown'}`,
        exp,
      });
    }
  },
});
