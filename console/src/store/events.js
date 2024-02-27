import { defineStore } from "pinia";
import { createPromiseClient } from "@connectrpc/connect";

import { EventsService } from "infinimesh-proto/build/es/eventbus/eventbus_connect";

import { useAppStore } from "@/store/app";
import { createConnectTransport } from "@connectrpc/connect-web";
import { EventKind } from "infinimesh-proto/build/es/eventbus/eventbus_pb";

export const useEventsStore = defineStore("events", () => {
  const appStore = useAppStore();
  const eventsApi = createPromiseClient(
    EventsService,
    createConnectTransport(appStore.transport_options)
  );

  const startEventsStream = async () => {
    for await (const event of eventsApi.subscribe()) {
      const { eventKind, ...data } = event.toJson();
      appStore.event_bus.publish(EventKind[eventKind], data);
    }
  };

  return { eventsApi, startEventsStream };
});
