<template>
  <n-collapse v-model:expanded-names="expanded">
    <n-collapse-item title="Reported State" name="reported">
      <template #header-extra v-if="reported">
        <n-button
          tertiary
          circle
          type="info"
          @click.stop.prevent="
            handleCopy(JSON.stringify(reported.data, null, 2))
          "
        >
          <template #icon>
            <n-icon><copy-outline /></n-icon>
          </template>
        </n-button>
      </template>
      <n-code
        language="json"
        :word-wrap="true"
        :code="
          reported
            ? JSON.stringify(reported.data, null, 2)
            : '// No State have been reported yet'
        "
      >
      </n-code>
    </n-collapse-item>
    <n-space
      justify="start"
      align="center"
      v-if="reported && expanded.includes('reported')"
    >
      <n-statistic label="Timestamp">
        <n-date-picker
          input-readonly
          :value="new Date(reported.timestamp).getTime()"
          type="datetime"
          disabled
          class="pseudo-disabled"
        />
      </n-statistic>
    </n-space>
    <n-collapse-item title="Desired State" name="desired">
      <template #header-extra v-if="desired">
        <n-button
          tertiary
          circle
          type="info"
          @click.stop.prevent="
            handleCopy(JSON.stringify(desired.data, null, 2))
          "
        >
          <template #icon>
            <n-icon><copy-outline /></n-icon>
          </template>
        </n-button>
      </template>
      <n-code
        language="json"
        :word-wrap="true"
        :code="
          desired
            ? JSON.stringify(desired.data, null, 2)
            : '// No Desired state have been set yet'
        "
      />
    </n-collapse-item>
    <n-space
      justify="start"
      align="center"
      v-if="desired && expanded.includes('desired')"
    >
      <n-statistic label="Timestamp">
        <n-date-picker
          input-readonly
          :value="new Date(desired.timestamp).getTime()"
          type="datetime"
          disabled
          class="pseudo-disabled"
        />
      </n-statistic>
    </n-space>
    <n-collapse-item title="Patch Desired" name="patch" v-if="patch">
      <template #header-extra v-if="validation == 'success'">
        <n-button
          tertiary
          round
          type="warning"
          @click.stop.prevent="handleSubmit"
        >
          Submit
        </n-button>
      </template>
      <n-input
        v-model:value="state"
        type="textarea"
        placeholder="Desired State"
        :status="validation"
      />
    </n-collapse-item>
  </n-collapse>
</template>

<script setup>
import { ref, computed } from "vue";
import {
  NCode,
  NCollapse,
  NCollapseItem,
  NDatePicker,
  useMessage,
  NInput,
  NSpace,
  NStatistic,
  NNumberAnimation,
  NButton,
  NIcon,
} from "naive-ui";
import { CopyOutline } from "@vicons/ionicons5";

const props = defineProps({
  state: {
    type: Object,
    required: true,
  },
  patch: {
    type: Boolean,
    default: false,
  },
});
const emit = defineEmits(["submit"]);

const expanded_val = ref(["reported"]);
const expanded = computed({
  get() {
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

const state = ref(JSON.stringify(desired.value.data, null, 2));
const validation = computed(() => {
  try {
    let d = JSON.parse(state.value);
    if (typeof d != "object") return "error";
    return "success";
  } catch {
    return "error";
  }
});

function handleSubmit() {
  emit("submit", state.value);
}
</script>