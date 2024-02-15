import { defineStore } from "pinia";
import { createPromiseClient } from "@connectrpc/connect";
import { NamespacesService } from "infinimesh-proto/build/es/node/node_connect";
import {
  EmptyMessage,
  JoinRequest,
} from "infinimesh-proto/build/es/node/node_pb";
import { Namespace } from "infinimesh-proto/build/es/node/namespaces/namespaces_pb";
import { useAppStore } from "@/store/app";
import { createConnectTransport } from "@connectrpc/connect-web";

const appStore = useAppStore();

export const useNSStore = defineStore("namespaces", {
  state: () => ({
    loading: false,
    selected: "",
    namespaces: {},
    namespacesApi: createPromiseClient(
      NamespacesService,
      createConnectTransport(appStore.transport_options)
    ),
  }),

  getters: {
    namespaces_list: (state) => {
      return Object.values(state.namespaces);
    },
  },

  actions: {
    async fetchNamespaces(no_cache = false) {
      this.loading = true;

      const data = await this.namespacesApi.list(new EmptyMessage());

      if (no_cache) {
        this.namespaces = data.namespaces.reduce((result, namespace) => {
          result[namespace.uuid] = namespace;

          return result;
        }, {});
      } else {
        this.namespaces = {
          ...this.namespaces,
          ...data.namespaces.reduce((result, namespace) => {
            result[namespace.uuid] = namespace;

            return result;
          }, {}),
        };
      }

      this.loading = false;
    },
    loadJoins(uuid) {
      return this.namespacesApi.joins(new Namespace(this.namespaces[uuid]));
    },
    join(namespace, account, access) {
      return this.namespacesApi.join(
        new JoinRequest({ namespace, account, access })
      );
    },
    create(namespace) {
      return this.namespacesApi.create(new Namespace(namespace));
    },

    /**
     * Updates a namespace.
     * @param {Namespace} namespace - The namespace object to be updated.
     * @returns {Promise<Namespace>} - A promise that resolves with the updated namespace.
     */
    update(namespace) {
      if (!namespace.config) namespace.config = {}
      const result = Namespace.fromJson(namespace)

      result.config = result.config.fromJson(namespace.config)
      return this.namespacesApi.update(result);
    },
    deletables(uuid) {
      return this.namespacesApi.deletables(
        new Namespace(this.namespaces[uuid])
      );
    },
    delete(uuid) {
      const namespace = this.namespaces[uuid];

      delete this.namespaces[uuid];
      return this.namespacesApi.delete(new Namespace(namespace));
    },
  },

  persist: {
    enabled: true,
    strategies: [{ storage: localStorage, key: "infinimesh.ns" }],
  },
});
