<template>
  <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
    :value="selected" />
</template>

<script setup>
import { ref, h, computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { NMenu } from "naive-ui";

import { renderIcon } from "@/utils";
import { GridOutline, HardwareChipOutline, PeopleOutline, GitNetworkOutline, ImagesOutline, ExtensionPuzzleOutline } from "@vicons/ionicons5";

import { useAppStore } from "@/store/app"
import { storeToRefs } from "pinia"

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

const { console_services } = storeToRefs(useAppStore())

const services = computed(() => {
  let r = []
  if (console_services.value.http_fs) {
    r.push({
      label: renderLabelLink("Media"),
      key: "Media",
      icon: renderIcon(ImagesOutline),
    })
  }
  return r
})

const menuOptions = ref([
  {
    label: renderLabelLink("DashboardMain", "Dashboard"),
    key: "Dashboard",
    icon: renderIcon(GridOutline),
  },
  {
    label: renderLabelLink("Devices"),
    key: "Devices",
    icon: renderIcon(HardwareChipOutline),
  },
  {
    label: renderLabelLink("Accounts"),
    key: "Accounts",
    icon: renderIcon(PeopleOutline),
  },
  {
    label: renderLabelLink("Namespaces"),
    key: "Namespaces",
    icon: renderIcon(GitNetworkOutline),
  },
  ...services.value,
  {
    label: renderLabelLink("Plugins", "Apps & Plugins"),
    key: "Plugins",
    icon: renderIcon(ExtensionPuzzleOutline),
  }
]);

const collapsed = computed(() => props.collapsed);
</script>