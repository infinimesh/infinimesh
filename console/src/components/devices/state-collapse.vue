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
      <n-code language="json" :word-wrap="true" :code="
        reported
          ? JSON.stringify(reported.data, null, 2)
          : '// No State have been reported yet'
      ">
      </n-code>
    </n-collapse-item>
    <n-grid responsive="screen" :collapsed-rows="2" v-if="reported && expanded.includes('reported')">
      <n-grid-item :span="8">
        <span>Timestamp</span>
      </n-grid-item>
      <n-grid-item :span="16">
        <n-date-picker input-readonly :value="new Date(reported.timestamp).getTime()" type="datetime" disabled
          class="pseudo-disabled" />
      </n-grid-item>
    </n-grid>
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
      <n-code language="json" :word-wrap="true" :code="
        desired
          ? JSON.stringify(desired.data, null, 2)
          : '// No Desired state have been set yet'
      " />
    </n-collapse-item>
    <n-grid responsive="screen" :collapsed-rows="2" v-if="desired && expanded.includes('desired')">
      <n-grid-item :span="8">
        <span>Timestamp</span>
      </n-grid-item>
      <n-grid-item :span="16">
        <n-date-picker input-readonly :value="new Date(desired.timestamp).getTime()" type="datetime" disabled
          class="pseudo-disabled" />
      </n-grid-item>
    </n-grid>
    <n-collapse-item title="Patch Desired" name="patch" v-if="patch">
      <template #header-extra v-if="validation == 'success'">
        <n-button tertiary round type="warning" @click.stop.prevent="handleSubmit">
          Submit
        </n-button>
      </template>
      <n-input v-model:value="desired_state" type="textarea" placeholder="Desired State" :status="validation" />
    </n-collapse-item>
  </n-collapse>
</template>

<script setup>
import { ref, computed, watch } from "vue";
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

const desired_state = ref(JSON.stringify(desired.value.data, null, 2));
watch(props.patch, () => {
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
</script>