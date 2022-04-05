<template>
  <n-collapse v-model:expanded-names="expanded">
    <n-collapse-item title="Reported State" name="reported">
      <n-code :code="reported ? JSON.stringify(reported.data, null, 2) : '// No State have been reported yet'" language="json" />
    </n-collapse-item>
    <n-space justify="space-between" align="center" v-if="reported && expanded.includes('reported')">
      <n-statistic label="Version">
        <n-number-animation
          :from="0"
          :to="reported.version"
          :active="true"
        />
      </n-statistic>
      <n-statistic label="Timestamp">
        <n-date-picker input-readonly :value="(new Date(reported.timestamp)).getTime()" type="datetime" disabled class="pseudo-disabled" />
      </n-statistic>
    </n-space>
    <n-collapse-item title="Desired State" name="desired">
      <n-code :code="desired ? JSON.stringify(desired.data, null, 2) : '// No Desired state have been set yet'" language="json" />
    </n-collapse-item>
    <n-space justify="space-between" align="center" v-if="desired && expanded.includes('desired')">
      <n-statistic label="Version">
        <n-number-animation
          :from="0"
          :to="desired.version"
          :active="true"
        />
      </n-statistic>
      <n-statistic label="Timestamp">
        <n-date-picker input-readonly :value="(new Date(desired.timestamp)).getTime()" type="datetime" disabled class="pseudo-disabled" />
      </n-statistic>
    </n-space>
  </n-collapse>
</template>

<script setup>
import { ref, computed } from "vue";
import {
  NCode, NCollapse, NCollapseItem, NDatePicker,
  NSpace, NStatistic, NNumberAnimation } from "naive-ui"


const props = defineProps({
  state: {
    type: Object,
    required: true,
  },
})

const expanded = ref(['reported'])

const reported = computed(() => {
  let state = props.state
  if (!state || !state.reported || state.reported.version == '0') {
    return false
  }
  return state.reported
})

const desired = computed(() => {
  let state = props.state
  if (!state || !state.desired || state.desired.version == '0') {
    return false
  }
  return state.desired
})

</script>