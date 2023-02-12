<template>
  <n-spin :show="patching" @mouseover="hover = true" @mouseleave="hover = false">
    <n-card hoverable :header-style="{ fontFamily: 'Exo 2' }" style="border-radius: 0">
      <template #header>
        <n-space aligh="center">
          {{ device.title }}
          <div :style="{ visibility: hover ? '' : 'hidden' }">
            <edit-dev-title-modal :device="device" @save="handleUpdateTitle" />
          </div>
        </n-space>
      </template>
      <template #header-extra>
        <n-tooltip trigger="hover" @click="handleUUIDClicked">
          <template #trigger>
            <n-tag :color="{ textColor: bulb_color, borderColor: bulb_color }" size="large" round
              @click="handleUUIDClicked">
              {{ device.uuid_short }}
            </n-tag>
          </template>
          {{ device.uuid }}
        </n-tooltip>
        <n-tooltip trigger="hover" @click="handleToggle">
          <template #trigger>
            <n-icon size="2vh" :color="bulb_color" style="margin-left: 1vw; cursor: pointer;"
              :class="toggle_animation ? 'jump-shaking-animation' : ''" @click="handleToggle">
              <bulb />
            </n-icon>
          </template>
          Click to toggle device(enable/disable)
        </n-tooltip>
      </template>

      <template #footer>
        <template v-if="show_ns">
          Namespace: <strong>{{ nss.namespaces[device.access.namespace]?.title || device.access.namespace }}</strong>
        </template><br />
        <template v-if="device.tags.length > 0">
          Tags:
          <n-tag type="warning" round v-for="tag in device.tags" :key="tag" style="margin-right: 3px">
            {{ tag }}
          </n-tag>
        </template>
        <n-space align="center" :style="{ visibility: hover ? '' : 'hidden', marginTop: '1rem' }">
          <edit-tags-modal :device="device" @save="handleUpdateTags" />
          <move v-if="access_lvl_conv(device) >= 3" type="device" :obj="device" @move="handleMove" />
        </n-space >
      </template>

      <template #action>
        <template v-if="plugin && plugin.kind == 'DEVICE'">
          <n-tabs type="segment" @update:value="handleStateTabChanged" :value="state_tab">
            <n-tab-pane :name="plugin.uuid" :tab="plugin.title">
              <div v-if="frame_url" style="width: 100%; height: max-content; overflow: visible;">
                <iframe :style="{ border: 'none', width: '100%', height: plugins.height(device.uuid) }" :src="frame_url"
                  ref="frame" scrolling="no" @load="iframeLoad"></iframe>
              </div>
              <n-alert title="Loading..." type="info" v-else>
                We're loading the device plugin frame
              </n-alert>
            </n-tab-pane>

            <n-tab-pane name="default" tab="JSON">
              <device-state-collapse :state="state" :patch="patch" :debug="debug" @submit="handlePatchDesired"
                @submit-debug="handlePatchReported" />
            </n-tab-pane>
          </n-tabs>

        </template>
        <template v-else>
          <device-state-collapse :state="state" :patch="patch" :debug="debug" @submit="handlePatchDesired"
            @submit-debug="handlePatchReported" />
        </template>

        <n-space justify="start" align="center" style="margin-top: 1vh">
          <n-button type="success" round tertiary :disabled="subscribed" @click="handleSubscribe">
            {{ subscribed? "Subscribed": "Subscribe" }}
          </n-button>

          <n-button v-if="access_lvl_conv(device) > 1" type="warning" round tertiary @click="patch = !patch">
            {{ patch? "Cancel Patch": "Patch Desired" }}
          </n-button>

          <n-button type="info" round tertiary @click="handleMakeToken">Make Device Token</n-button>

          <n-popconfirm @positive-click="handleDelete" v-if="access_lvl_conv(device) > 2">
            <template #trigger>
              <n-button type="error" round secondary>Delete</n-button>
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

      <status-corner :connection="state.connection" />
    </n-card>
    <n-modal :show="plugin_edit_modal && plugin_edit_modal.show" preset="dialog" size="huge"
      @update:show="(v) => !v && (patch = false)" style="width: 90vw">
      <n-space justify="space-around">
        <iframe :style="{ border: 'none', width: '85vw', height: '80vh' }" :src="plugin_edit_modal.frame">
        </iframe>
      </n-space>
    </n-modal>
  </n-spin>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, watch } from "vue";
import {
  NCard, NTooltip, NAlert,
  NIcon, useMessage, NModal,
  NSpin, useLoadingBar,
  NTag, NSpace, NButton,
  NPopconfirm, NTabs, NTabPane
} from "naive-ui";

import { useDevicesStore } from "@/store/devices";
import { useNSStore } from "@/store/namespaces";
import { baseURL, useAppStore } from "@/store/app";
import { usePluginsStore } from "@/store/plugins";

import { access_lvl_conv } from "@/utils/access";
import { storeToRefs } from "pinia";

const Bulb = defineAsyncComponent(() => import("@vicons/ionicons5/Bulb"))
const BugOutline = defineAsyncComponent(() => import("@vicons/ionicons5/BugOutline"))

const EditDevTitleModal = defineAsyncComponent(() => import('./edit-dev-title-modal.vue'))
const EditTagsModal = defineAsyncComponent(() => import("./edit-tags-modal.vue"))
const DeviceStateCollapse = defineAsyncComponent(() => import("./state-collapse.vue"))

const StatusCorner = defineAsyncComponent(() => import("./status-corner.vue"))

const Move = defineAsyncComponent(() => import("@/components/namespaces/move.vue"))

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

const state = computed(() => store.device_state(device.value.uuid))

const bulb_color = computed(() => {
  if (device.value.enabled === null) return "gray";
  return device.value.enabled ? "#52c41a" : "#eb2f96";
});
const toggle_animation = ref(false)

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

const { dev, theme } = storeToRefs(useAppStore())

const plugins = usePluginsStore()
const { current: plugin } = storeToRefs(plugins)
let token = false
const frame = ref(null)
const frame_url = ref(false)

async function frame_src(view = 'viewUrl') {
  if (!plugin.value || plugin.value.kind != 'DEVICE') {
    return
  }

  const uuid = props.device.uuid

  if (!token) {
    token = await store.makeDevicesToken([uuid], true)
  }

  let params = {
    token, device: uuid,
    theme: theme.value, api: baseURL,
    vars: plugin.value.vars,
    ...state.value
  }

  const src = `${plugin.value.deviceConf[view]}?a=${btoa(JSON.stringify(params))}`
  return src
}

async function prepare_frame() {
  frame_url.value = false

  frame_url.value = await frame_src()
}
watch([plugin, theme, state], prepare_frame)
prepare_frame()

const bar = useLoadingBar();
const patch = ref(false);
const patching = ref(false);
async function handlePatchDesired(state) {
  patching.value = true;
  await store.patchDesiredState(device.value.uuid, state, bar);
  patch.value = false;
  patching.value = false;
}

const state_tab = ref((plugin.value && plugin.value.uuid) || 'default')
const plugin_edit_modal = ref(false)

function handleStateTabChanged(v) {
  state_tab.value = v
}

watch(patch, async (n) => {
  if (!plugin.value || plugin.value.kind != 'DEVICE') return

  if (!n) {
    plugin_edit_modal.value = false
    return
  }

  plugin_edit_modal.value = {
    show: true,
    frame: await frame_src('desiredUrl')
  }
})


const debug = ref(false) // Enables Reported state editor

// API handlers
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
  toggle_animation.value = false
  await store.toggle(device.value.uuid, bar);
  toggle_animation.value = true
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
    await store.updateDevice(device.value.uuid, { title: device.value.title, tags: tags })
    resolve()
  } catch (e) {
    reject(e)
  }
}

async function handleUpdateTitle(title, resolve, reject) {
  try {
    await store.updateDevice(device.value.uuid, { title: title, tags: device.value.tags })
    resolve()
  } catch (e) {
    reject(e)
  }
}

async function handleMove(ns, resolve, reject) {
  try {
    await store.moveDevice(device.value.uuid, ns)
    resolve()
  } catch (e) {
    reject(e)
  }
}

const nss = useNSStore()

const hover = ref(false)
</script>