import { watchEffect, onBeforeMount, computed } from "vue";
import { useNSStore } from "@/store/namespaces";

export default function (name = "base", options, changeOptions) {
  const nsStore = useNSStore();

  const currentNamespace = computed(() => nsStore.selected);

  watchEffect(
    () => {
      const paginationState = JSON.parse(
        localStorage.getItem("pagination") || `{}`
      );
      const currentState = paginationState[name] || {};
      currentState[currentNamespace.value] = { limit: options.limit.value };
      paginationState[name] = currentState;

      localStorage.setItem("pagination", JSON.stringify(paginationState));
    },
    { flush: "post" }
  );

  watchEffect(
    () => {
      const paginationState = JSON.parse(
        localStorage.getItem("pagination") || `{}`
      );
      const currentState = paginationState[name] || {};
      if (currentState[currentNamespace.value]?.limit) {
        changeOptions({ limit: currentState[currentNamespace.value].limit });
      }
    },
    { flush: "pre" }
  );
}
