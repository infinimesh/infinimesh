import { defineStore } from "pinia";
import { useAppStore } from "@/store/app";

const as = useAppStore();

export const useSessionsStore = defineStore("sessions", () => {

    async function get() {
        const { data } = await as.http.get("/sessions");
        return data;
    }

    async function activity() {
        const { data } = await as.http.get('/sessions/activity');
        return data;
    }

    async function revoke(sid) {
        await as.http.delete(`/sessions/${sid}`);
    }
    return {
        get, activity, revoke
    }
})