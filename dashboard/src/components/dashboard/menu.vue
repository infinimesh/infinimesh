<template>
  <n-menu
    :collapsed="collapsed"
    :collapsed-width="64"
    :collapsed-icon-size="22"
    :options="menuOptions"
    :value="selected"
  />
</template>

<script setup>
import { ref, h, computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { NMenu } from "naive-ui";

import { renderIcon } from "@/utils";
import { HardwareChipOutline } from "@vicons/ionicons5";

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false,
  },
});

const route = useRoute();
const selected = computed(() => route.name);
const menuOptions = ref([
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "Devices",
          },
        },
        { default: () => "Devices" }
      ),
    key: "Devices",
    icon: renderIcon(HardwareChipOutline),
  },
]);

const collapsed = computed(() => props.collapsed);
</script>