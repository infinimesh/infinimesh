<template>
  <n-spin :show="patching" @mouseover="hover = true" @mouseleave="hover = false">
    <n-card hoverable :title="device.title" :header-style="{ fontFamily: 'Exo 2' }" style="border-radius: 0">
      <template #header-extra>
        <n-tooltip trigger="hover" @click="handleUUIDClicked">
          <template #trigger>
            <n-tag :color="{ textColor: bulb_color, borderColor: bulb_color }" size="large" round
              @click="handleUUIDClicked">
              {{  device.uuid_short  }}
            </n-tag>
          </template>
          {{  device.uuid  }}
        </n-tooltip>
        <n-tooltip trigger="hover" @click="handleToggle">
          <template #trigger>
            <n-icon size="2vh" :color="bulb_color" style="margin-left: 1vw" @click="handleToggle">
              <bulb />
            </n-icon>
          </template>
          Click to toggle device(enable/disable)
        </n-tooltip>
      </template>

      <template #footer>
        <template v-if="show_ns">
          Namespace: <strong>{{  nss.namespaces[device.access.namespace]?.title || device.access.namespace  }}</strong>
        </template><br />
        <template v-if="device.tags.length > 0">
          Tags:
          <n-tag type="warning" round v-for="tag in device.tags" :key="tag" style="margin-right: 3px">
            {{  tag  }}
          </n-tag>
        </template>
        <div :style="{ visibility: hover ? '' : 'hidden', marginTop: '1rem' }">
          <edit-tags-modal :device="device" @save="handleUpdateTags" />
        </div>
      </template>

      <template #action>
        <device-state-collapse :state="store.device_state(device.uuid)" :patch="patch" :debug="debug"
          @submit="handlePatchDesired" @submit-debug="handlePatchReported" />
        <n-space justify="start" align="center" style="margin-top: 1vh">
          <n-button type="success" round tertiary :disabled="subscribed" @click="handleSubscribe">
            {{  subscribed ? "Subscribed" : "Subscribe"  }}
          </n-button>

          <n-button v-if="access_lvl_conv(device) > 1" type="warning" round tertiary @click="patch = !patch">
            {{  patch ? "Cancel Patch" : "Patch Desired"  }}
          </n-button>

          <n-button type="info" round tertiary @click="handleMakeToken">Make Device Token</n-button>

          <n-popconfirm @positive-click="handleDelete">
            <template #trigger>
              <n-button v-if="access_lvl_conv(device) > 2" type="error" round secondary>Delete</n-button>
            </template>
            Are you sure about deleting this device?
          </n-popconfirm>

          <template v-if="dev">

            <n-button round tertiary type="success" @click="debug = false" v-if="debug">
              <template #icon>
                <n-icon>
                  <bug-outline />
                </n-icon>
              </template>
              Cancel debug
            </n-button>
            <n-button secondary circle type="success" @click="debug = true" v-else>
              <template #icon>
                <n-icon>
                  <bug-outline />
                </n-icon>
              </template>
            </n-button>

          </template>
        </n-space>
      </template>
    </n-card>
  </n-spin>
</template>

<script setup>
import { ref, computed, defineAsyncComponent } from "vue";
import {
  NCard,
  NTooltip,
  NIcon,
  useMessage,
  NSpin,
  useLoadingBar,
  NTag,
  NSpace,
  NButton,
  NPopconfirm,
} from "naive-ui";

import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { useAppStore } from "@/store/app";

import { access_lvl_conv } from "@/utils/access";
import { storeToRefs } from "pinia";

const Bulb = defineAsyncComponent(() => import("@vicons/ionicons5/Bulb"))
const BugOutline = defineAsyncComponent(() => import("@vicons/ionicons5/BugOutline"))

const EditTagsModal = defineAsyncComponent(() => import("./edit-tags-modal.vue"))
const DeviceStateCollapse = defineAsyncComponent(() => import("./state-collapse.vue"))

const props = defineProps({
  device: {
    type: Object,
    required: true,
  },
  show_ns: {
    type: Boolean,
    default: false,
  },
});

const device = computed(() => {
  let r = props.device;
  r.uuid_short = r.uuid.substr(0, 8);
  return r;
});

const store = useDevicesStore();

const subscribed = computed(() => {
  return store.device_subscribed(device.value.uuid);
});

const bulb_color = computed(() => {
  if (device.value.enabled === null) return "gray";
  return device.value.enabled ? "#52c41a" : "#eb2f96";
});

const message = useMessage();
async function handleUUIDClicked() {
  try {
    await navigator.clipboard.writeText(device.value.uuid);
    message.success("Device UUID copied to clipboard");
  } catch {
    message.error("Failed to copy device UUID to clipboard");
  }
}

function handleSubscribe() {
  store.subscribe([device.value.uuid]);
}

const { dev } = storeToRefs(useAppStore())

const bar = useLoadingBar();
const patch = ref(false);
const patching = ref(false);
async function handlePatchDesired(state) {
  patching.value = true;
  await store.patchDesiredState(device.value.uuid, state, bar);
  patch.value = false;
  patching.value = false;
}

const debug = ref(false)
async function handlePatchReported(state) {
  patching.value = true;
  await store.patchReportedState(device.value.uuid, state, bar);
  patch.value = false;
  patching.value = false;
}

async function handleDelete() {
  patching.value = true;
  await store.deleteDevice(device.value.uuid, bar);
  patching.value = false;
}

async function handleToggle() {
  store.toggle(device.value.uuid, bar);
}

async function handleMakeToken() {
  let token = await store.makeDevicesToken(
    [device.value.uuid],
    access_lvl_conv(device.value) > 1
  );
  try {
    await navigator.clipboard.writeText(token);
    message.success("Device Token copied to clipboard");
  } catch {
    message.error("Failed to copy device token to clipboard");
  }
}

async function handleUpdateTags(tags, resolve, reject) {
  try {
    console.log('updating', device.value.uuid, tags)
    await store.updateDevice(device.value.uuid, { title: device.value.title, tags: tags })
    resolve()
  } catch (e) {
    reject(e)
  }
}

const nss = useNSStore()

const hover = ref(false)
</script>