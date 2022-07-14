import { useAppStore } from "@/store/app";
import { defineStore } from "pinia";

const as = useAppStore();

export const useIStore = defineStore("internal", {
    state: () => ({
        ldap_providers: {}
    }),

    actions: {
        async getLDAPProviders() {
            try {
                const { data } = await as.http.get("/i/ldapp");
                this.ldap_providers = data.providers
            } catch (e) {
                console.warn("Error while getting LDAP providers", e)
            }
        }
    }
})