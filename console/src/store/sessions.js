import { defineStore } from "pinia";
import { createPromiseClient } from "@connectrpc/connect";

import { SessionsService } from "infinimesh-proto/build/es/node/node_connect"
import { Session } from "infinimesh-proto/build/es/node/sessions/sessions_pb";
import { EmptyMessage } from "infinimesh-proto/build/es/node/node_pb";

import { useAppStore } from "@/store/app";

export const useSessionsStore = defineStore("sessions", () => {
  const appStore = useAppStore();
  const sessionsApi = createPromiseClient(SessionsService, appStore.transport);

  async function get() {
    return await sessionsApi.get(new EmptyMessage());
  }

  async function activity() {
    return await sessionsApi.getActivity(new EmptyMessage());
  }

  async function revoke(session) {
    return await sessionsApi.revoke(new Session(session));
  }

  return { get, activity, revoke };
})