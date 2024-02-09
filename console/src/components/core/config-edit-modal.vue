<template>
  <modal-button
    type="info"
    min-width="30vw"
    submit-text="Submit"
    cancel-text="Cancel"
    @submit="handleSubmit"
    @show="reset"
  >
    <template #icon> <cog-outline /> </template>
    <template #header>
      Manage Configuration Data for <strong>{{ o.title }}</strong>
    </template>

    <n-space vertical justify="space-between">
      <n-input
        v-model:value="config"
        type="textarea"
        placeholder="Desired State"
        :status="validation"
      />
    </n-space>
  </modal-button>
</template>

<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import { NSpace, NInput } from 'naive-ui'

const CogOutline = defineAsyncComponent(
  () => import('@vicons/ionicons5/CogOutline')
)
const ModalButton = defineAsyncComponent(
  () => import('@/components/core/modal-button.vue')
)

const props = defineProps({
  o: { type: Object, required: true }
})
const emits = defineEmits(['submit'])

const item = ref(JSON.parse(JSON.stringify(props.o)))
const config = ref(JSON.stringify(props.o.config ?? {}, null, 2))

const validation = computed(() => {
  try {
    const result = JSON.parse(config.value)

    if (typeof result !== 'object') return 'error'
    return 'success'
  } catch {
    return 'error'
  }
})

function handleSubmit(close) {
  if (validation.value) {
    item.value.config = JSON.parse(config.value)

    emits('submit', item.value)
    close()
  }
}

function reset() {
  config.value = JSON.stringify(props.o.config ?? {}, null, 2)
}

</script>