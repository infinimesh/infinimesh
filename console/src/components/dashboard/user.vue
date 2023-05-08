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
import { defineAsyncComponent, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { NButton, NDropdown, NIcon } from "naive-ui";
import { renderIcon } from "@/utils";
import { useAppStore } from "@/store/app";
import { useRouter } from "vue-router";

const Person = defineAsyncComponent(() => import("@vicons/ionicons5/Person"));
const LogOutOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LogOutOutline"));
const LockClosedOutline = defineAsyncComponent(() => import("@vicons/ionicons5/LockClosedOutline"));
const CodeSlashOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CodeSlashOutline"));
const KeyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/KeyOutline"));
const CogOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CogOutline"));

const router = useRouter();
const store = useAppStore();

const { me, dev } = storeToRefs(store);

const options = ref([
  {
    key: "credentials",
    label: "Manage Credentials",
    icon: renderIcon(LockClosedOutline),
    props: {
      onClick: () => {
        router.push({ name: "Credentials" });
      }
    },
  },
  {
    key: 'token',
    label: 'Personal Access Token',
    icon: renderIcon(KeyOutline),
    props: {
      onClick: () => {
        router.push({ name: "Tokens" });
      }
    }
  },
  {
    type: 'divider',
    key: 'd1'
  },
  {
    key: "settings",
    label: "Settings",
    icon: renderIcon(CogOutline),
    props: {
      onClick: () => {
        router.push({ name: "Profile" });
      },
    },
  },
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
  }
]);

function addNoDev() {
  if (dev.value) {
    options.value.splice(1, 0, {
      key: "nodev",
      label: "Turn off Develeper Mode",
      icon: renderIcon(CodeSlashOutline),
      props: {
        onClick: () => {
          dev.value = false
          options.value.splice(1, 1)
        }
      }
    })
  }
}
watch(dev, addNoDev)
addNoDev()
</script>