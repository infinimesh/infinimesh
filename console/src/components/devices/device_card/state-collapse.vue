<template>
  <n-collapse v-model:expanded-names="expanded">
    <n-collapse-item title="Reported State" name="reported">
      <template #header-extra v-if="reported">
        <n-button tertiary circle type="info" @click.stop.prevent="
          handleCopy(JSON.stringify(reported.data, null, 2))
        ">
          <template #icon>
            <n-icon>
              <copy-outline />
            </n-icon>
          </template>
        </n-button>
      </template>

      <n-grid responsive="screen" :collapsed-rows="2" v-if="reported">
        <n-grid-item :span="8">
          <span>Timestamp</span>
        </n-grid-item>
        <n-grid-item :span="16">
          <n-date-picker input-readonly :value="timestamppb_to_time(reported.timestamp)" type="datetime" disabled
            class="pseudo-disabled" />
        </n-grid-item>
      </n-grid>

      <div class="reported-state" :style="{ maxHeight: calculatedHeight }">
        <n-code language="json" :word-wrap="true" :code="
          reported
            ? JSON.stringify(reported.data, null, 2)
            : '// No State have been reported yet'
        ">
        </n-code>
      </div>
    </n-collapse-item>

    <n-collapse-item title="Patch Reported" name="reported_patch" v-if="debug">
      <template #header-extra>
        <n-button tertiary round type="warning" @click.stop.prevent="handleSubmitReported"
          :disabled="reported_validation != 'success'">
          Submit
        </n-button>
      </template>
      <n-input v-model:value="reported_state" type="textarea" placeholder="Reported State"
        :status="reported_validation" />
    </n-collapse-item>

    <n-collapse-item title="Desired State" name="desired">
      <template #header-extra v-if="desired">
        <n-button tertiary circle type="info" @click.stop.prevent="
          handleCopy(JSON.stringify(desired.data, null, 2))
        ">
          <template #icon>
            <n-icon>
              <copy-outline />
            </n-icon>
          </template>
        </n-button>
      </template>

      <n-grid responsive="screen" :collapsed-rows="2" v-if="desired">
        <n-grid-item :span="8">
          <span>Timestamp</span>
        </n-grid-item>
        <n-grid-item :span="16">
          <n-date-picker input-readonly :value="timestamppb_to_time(desired.timestamp)" type="datetime" disabled
            class="pseudo-disabled" />
        </n-grid-item>
      </n-grid>

      <div class="reported-state" :style="{ maxHeight: calculatedHeight }">
        <n-code language="json" :word-wrap="true" :code="
          desired
            ? JSON.stringify(desired.data, null, 2)
            : '// No Desired state have been set yet'
        " />
      </div>
    </n-collapse-item>

    <n-collapse-item title="Patch Desired" name="patch" v-if="patch">
      <template #header-extra>
        <n-button tertiary round type="warning" @click.stop.prevent="handleSubmit" :disabled="validation != 'success'">
          Submit
        </n-button>
      </template>
      <n-input v-model:value="desired_state" type="textarea" placeholder="Desired State" :status="validation" />
    </n-collapse-item>
  </n-collapse>
</template>

<script setup>
import { ref, computed, watch, toRef, defineAsyncComponent } from "vue";
import {
  NCode,
  NCollapse,
  NCollapseItem,
  NDatePicker,
  useMessage,
  NInput,
  NButton,
  NIcon,
  NGrid,
  NGridItem,
} from "naive-ui";

import { Shadow } from "infinimesh-proto/build/es/shadow/shadow_pb";
import { Timestamp } from "@bufbuild/protobuf";

const CopyOutline = defineAsyncComponent(() => import("@vicons/ionicons5/CopyOutline"))

const props = defineProps({
  state: {
    type: Object,
    required: true,
  },
  patch: {
    type: Boolean,
    default: false,
  },
  debug: {
    type: Boolean,
    default: false,
  }
});
const emit = defineEmits(["submit"]);

const expanded_val = ref(["reported"]);
const expanded = computed({
  get() {
    if (props.debug) {
      return ["reported", "reported_patch", "desired", "patch"];
    }
    if (props.patch) {
      return ["reported", "desired", "patch"];
    }
    return expanded_val.value;
  },
  set(val) {
    expanded_val.value = val;
  },
});

const reported = computed(() => {
  let state = props.state;
  if (!state || !state.reported || !state.reported.timestamp) {
    return false;
  }
  return state.reported;
});

const desired = computed(() => {
  let state = props.state;
  if (!state || !state.desired || !state.desired.timestamp) {
    return false;
  }
  return state.desired;
});

const message = useMessage();
async function handleCopy(state) {
  try {
    await navigator.clipboard.writeText(state);
    message.success("Device State copied to clipboard");
  } catch {
    message.error("Failed to copy device State to clipboard");
  }
}

const desired_state = ref(JSON.stringify(desired.value.data, null, 2));
watch(toRef(props.patch), () => {
  desired_state.value = ref(JSON.stringify(desired.value.data, null, 2));
})

const validation = computed(() => {
  try {
    let d = JSON.parse(desired_state.value);
    if (typeof d != "object") return "error";
    return "success";
  } catch {
    return "error";
  }
});

function handleSubmit() {
  emit("submit", JSON.parse(desired_state.value));
}

const reported_state = ref(JSON.stringify(reported.value.data, null, 2));
watch(toRef(props.debug), () => {
  reported_state.value = ref(JSON.stringify(reported.value.data, null, 2));
})

const reported_validation = computed(() => {
  try {
    let d = JSON.parse(reported_state.value);
    if (typeof d != "object") return "error";
    return "success";
  } catch {
    return "error";
  }
});

const calculatedHeight = computed(() => {
  const maxRows = STATE_MAX_ROWS || 10;
  const lineHeight = parseFloat(window.getComputedStyle(document.body).lineHeight);
  return (parseInt(maxRows) * lineHeight) + 'px';
});

function handleSubmitReported() {
  emit("submit-debug", JSON.parse(reported_state.value));
}

/*
  * Convert a timestamppb to a Date object
  * @param {Timestamp} timestamp
  * @returns {number}
  */
function timestamppb_to_time(timestamp) {
  return Number(timestamp.seconds) * 1000
}
</script>

<style scoped>
.reported-state {
  max-height: 224px;
  overflow-y: auto;
}
</style>