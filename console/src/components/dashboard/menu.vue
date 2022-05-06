<template>
  <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
    :value="selected" />
</template>

<script setup>
import { ref, h, computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { NMenu } from "naive-ui";

import { renderIcon } from "@/utils";
import { HardwareChipOutline, PeopleOutline, GitNetworkOutline } from "@vicons/ionicons5";

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false,
  },
});

const route = useRoute();
const selected = computed(() => route.name);

function renderLabelLink(route, label = false) {
  if (!label) {
    label = route
  }
  return () => h(
    RouterLink,
    {
      to: {
        name: route,
      },
    },
    { default: () => label }
  )
}

const menuOptions = ref([
  {
    label: renderLabelLink("Accounts"),
    key: "Accounts",
    icon: renderIcon(PeopleOutline),
  },
  {
    label: renderLabelLink("Devices"),
    key: "Devices",
    icon: renderIcon(HardwareChipOutline),
  },
  {
    label: renderLabelLink("Namespaces"),
    key: "Namespaces",
    icon: renderIcon(GitNetworkOutline),
  }
]);

const collapsed = computed(() => props.collapsed);
</script>