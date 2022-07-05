import { useAppStore } from "@/store/app";
import { defineStore } from "pinia";

import { check_token_expired } from "@/utils/access";

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
                check_token_expired(e, as)
            }
        }
    }
})