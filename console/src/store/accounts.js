import { useAppStore } from "@/store/app";
import { useNSStore } from "@/store/namespaces";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

import { access_lvl_conv } from "@/utils/access";
import { createPromiseClient } from "@connectrpc/connect";
import { AccountsService } from "infinimesh-proto/build/es/node/node_connect";
import {
  Account,
  CreateRequest,
} from "infinimesh-proto/build/es/node/accounts/accounts_pb";
import {
  EmptyMessage,
  MoveRequest,
  SetCredentialsRequest,
  TokenRequest,
} from "infinimesh-proto/build/es/node/node_pb";
import { createConnectTransport } from "@connectrpc/connect-web";

import { transport } from "infinimesh-proto/mocks/es/accounts";

export const useAccountsStore = defineStore("accounts", () => {
  const as = useAppStore();
  const nss = useNSStore();

  const loading = ref(false);
  const accounts = ref({});

  const accounts_ns_filtered = computed(() => {
    let ns = nss.selected;
    let pool = Object.values(accounts.value)
      .map((acc) => {
        acc.sorter = acc.enabled + access_lvl_conv(acc);
        return acc;
      })
      .sort((a, b) => b.sorter - a.sorter);
    if (ns == "all") {
      return pool;
    }
    return pool.filter((a) => a.access.namespace == ns);
  });

  const accountsApi = computed(() => {
    return createPromiseClient(
      AccountsService,
      import.meta.env.VITE_MOCK
        ? transport
        : createConnectTransport(as.transport_options)
    );
  });

  async function sync_me() {
    const data = await accountsApi.value.get({ uuid: "me" });

    as.me = { ...as.me, ...data };
  }

  async function fetchAccounts(no_cache = false) {
    loading.value = true;

    const { accounts: pool } = await accountsApi.value.list(new EmptyMessage());

    if (no_cache) {
      accounts.value = pool.reduce((result, account) => {
        result[account.uuid] = account;

        return result;
      }, {});
    } else {
      accounts.value = {
        ...accounts.value,
        ...pool.reduce((result, account) => {
          result[account.uuid] = account;

          return result;
        }, {}),
      };
    }

    loading.value = false;
  }

  async function createAccount(request, bar) {
    bar.start();
    try {
      await accountsApi.value.create(new CreateRequest(request));

      fetchAccounts();
      bar.finish();
      return false;
    } catch (e) {
      console.error(e);
      bar.error();
      return e;
    }
  }

  async function updateAccount(account, bar) {
    if (bar) bar.start();
    try {
      if (!account.config) account.config = {};
      const result = new Account(account);

      result.config = result.config.fromJson(account.config);
      await accountsApi.value.update(result);

      fetchAccounts();
      if (bar) bar.finish();
      return false;
    } catch (e) {
      console.error(e);
      if (bar) bar.error();
      return e;
    }
  }
  function updateDefaultNamespace(account, ns) {
    let acc = accounts.value[account];
    return updateAccount({
      ...acc,
      defaultNamespace: ns,
    });
  }
  async function toggle(uuid, bar) {
    bar.start();

    try {
      await accountsApi.value.toggle(new Account(accounts.value[uuid]));

      fetchAccounts();
      bar.finish();
    } catch (e) {
      console.error(e);
      bar.error();
      return;
    }
  }
  function deletables(uuid) {
    return accountsApi.value.deletables(new Account(accounts.value[uuid]));
  }
  async function deleteAccount(uuid, bar) {
    bar.start();
    try {
      await accountsApi.value.delete(new Account(accounts.value[uuid]));

      fetchAccounts();
      bar.finish();
    } catch (e) {
      console.error(e);
      bar.error();
    }
  }
  async function moveAccount(uuid, namespace) {
    try {
      await accountsApi.value.move(new MoveRequest({ uuid, namespace }));

      accounts.value[uuid].access.namespace = namespace;
    } catch (err) {
      console.error(err);
      throw `Error Moving Device: ${err.message}`;
    }
  }
  function getCredentials(uuid) {
      return accountsApi.value.getCredentials({ uuid });
  }
  async function setCredentials(uuid, credentials, bar) {
    bar.start();
    try {
      await accountsApi.value.setCredentials(
        new SetCredentialsRequest({ uuid, credentials })
      );

      fetchAccounts();
      bar.finish();
    } catch (e) {
      console.error(e);
      bar.error();
    }
  }
  function token(data) {
    return accountsApi.value.token(new TokenRequest(data));
  }
  function tokenFor(account, exp = 0) {
    let res = {};
    try {
      res = new UAParser(navigator.userAgent).getResult();
    } catch (e) {
      console.warn("Failed to get user agent", e);
    }

    return token({
      uuid: account,
      client: `Console Admin | ${res.os?.name ?? "Unknown"} | ${
        res.browser?.name ?? "Unknown"
      }`,
      exp,
    });
  }

  return {
    loading,
    accounts,
    accounts_ns_filtered,
    accountsApi,
    sync_me,
    fetchAccounts,
    createAccount,
    updateAccount,
    updateDefaultNamespace,
    toggle,
    deletables,
    deleteAccount,
    moveAccount,
    getCredentials,
    setCredentials,
    token,
    tokenFor,
  };
});
