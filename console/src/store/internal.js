import { ref, computed } from "vue"
import { defineStore } from "pinia"
import { createPromiseClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import { InternalService } from "infinimesh-proto/build/es/node/node_connect"
import { EmptyMessage } from "infinimesh-proto/build/es/node/node_pb"
import { useAppStore } from "@/store/app.js"

export const useIStore = defineStore("internal", () => {
  const appStore = useAppStore()

  const internalApi = computed(() => {
    return createPromiseClient(
      InternalService,
      import.meta.env.VITE_MOCK
        ? {}
        : createConnectTransport(appStore.transport_options)
    )
  })
  const ldap_providers = ref({})

  return {
    async getLDAPProviders() {
      try {
        const { providers } = await internalApi.value.getLDAPProviders(new EmptyMessage())

        ldap_providers.value = providers
      } catch (error) {
        console.warn("Error while getting LDAP providers", error)
      }
    }
  }
})