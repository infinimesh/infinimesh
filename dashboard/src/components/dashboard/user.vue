<template>
  <n-dropdown :options="options">
    <n-button>
      <template #icon>
        <n-icon>
          <Person />
        </n-icon>
      </template>
      {{ me.title }}
    </n-button>
  </n-dropdown>
</template>

<script setup>
import { h, ref } from "vue";
import { storeToRefs } from "pinia";
import { NButton, NDropdown, NIcon } from "naive-ui";
import { renderIcon } from "@/utils";
import { Person, LogOutOutline } from "@vicons/ionicons5";
import { useAppStore } from "@/store/app";
import { useRouter } from "vue-router";

const router = useRouter();
const store = useAppStore();

const { me } = storeToRefs(store);

const options = ref([
  {
    key: "logout",
    label: "Logout",
    icon: renderIcon(LogOutOutline),
    props: {
      onClick: () => {
        store.logout();
        router.push({ name: "Login" });
      },
    },
  },
]);
</script>