import { ref, computed, watch, reactive } from 'vue'
import { defineStore } from 'pinia'
import { createPromiseClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'

import { NamespacesService } from 'infinimesh-proto/build/es/node/node_connect'
import { EmptyMessage, JoinRequest } from 'infinimesh-proto/build/es/node/node_pb'
import { Namespace } from 'infinimesh-proto/build/es/node/namespaces/namespaces_pb'
import { transport } from 'infinimesh-proto/mocks/es/namespaces'
import { useAppStore } from './app.js'

export const useNSStore = defineStore('namespaces', () => {
  const appStore = useAppStore()

  const loading = ref(false)
  const selected = ref('')
  const namespaces = ref({})

  const state = reactive({ loading, selected, namespaces })
  const rawState = JSON.parse(localStorage.getItem('infinimesh.ns'))

  watch(() => state, (value) => {
    localStorage.setItem('infinimesh.ns', JSON.stringify(value))
  }, { deep: true })

  if (rawState) {
    Object.keys(state).forEach((key) => {
      state[key] = rawState[key]
    })
  }

  const namespaces_list = computed(() =>
    Object.values(namespaces.value)
  )
  const namespacesApi = computed(() =>
    createPromiseClient(
      NamespacesService,
      (import.meta.env.VITE_MOCK) ? transport : createConnectTransport(appStore.transport_options)
    )
  )

  async function fetchNamespaces(no_cache = false) {
    loading.value = true

    const data = await namespacesApi.value.list(new EmptyMessage())

    if (no_cache) {
      namespaces.value = data.namespaces.reduce((result, namespace) => {
        result[namespace.uuid] = namespace

        return result
      }, {})
    } else {
      namespaces.value = {
        ...namespaces.value,
        ...data.namespaces.reduce((result, namespace) => {
          result[namespace.uuid] = namespace

          return result
        }, {}),
      }
    }

    loading.value = false
  }

  function loadJoins(uuid) {
    return namespacesApi.value.joins(new Namespace(namespaces.value[uuid]))
  }

  function join(namespace, account, access) {
    return namespacesApi.value.join(
      new JoinRequest({ namespace, account, access })
    )
  }

  function create(namespace) {
    return namespacesApi.value.create(new Namespace(namespace))
  }

  /**
   * Updates a namespace.
   * @param {Namespace} namespace - The namespace object to be updated.
   * @returns {Promise<Namespace>} - A promise that resolves with the updated namespace.
   */
  function update(namespace) {
    if (!namespace.config) namespace.config = {}
    const result = Namespace.fromJson(namespace)

    result.config = result.config.fromJson(namespace.config)
    return namespacesApi.value.update(result)
  }

  function deletables(uuid) {
    return namespacesApi.value.deletables(
      new Namespace(namespaces.value[uuid])
    )
  }

  function remove(uuid) {
    const namespace = namespaces.value[uuid]

    delete namespaces.value[uuid]
    return namespacesApi.value.delete(new Namespace(namespace))
  }

  return {
    loading,
    selected,
    namespaces,
    namespaces_list,
    namespacesApi,

    fetchNamespaces,
    loadJoins,
    join,
    create,
    update,
    deletables,
    remove
  }
})
