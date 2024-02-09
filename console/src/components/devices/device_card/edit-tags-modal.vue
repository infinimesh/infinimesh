<template>
  <n-button round tertiary type="info" @click="show = true">
    <template #icon>
      <n-icon> <add-outline /> </n-icon>
    </template>
    {{ (value.length > 0) ? 'Edit' : 'Add' }} {{ entity }}
  </n-button>

  <n-modal v-model:show="show">
    <n-spin :show="loading">
      <template #description>
        Updating {{ entity }}...
      </template>

      <n-card
        style="min-width: 30vw; max-width: 90vw"
        size="huge"
        role="dialog"
        aria-modal="true"
        :bordered="false"
      >
        <template #header>
          Edit Device({{ device.title }}) {{ entity }}
        </template>

        <template #header-extra>
          <n-button quaternary circle size="large" @click="show = false">
            <template #icon>
              <n-icon> <close-outline /> </n-icon>
            </template>
          </n-button>
        </template>

        <n-space vertical justify="space-between">
          <n-dynamic-tags v-model:value="tags" round type="warning" size="large">
            <template v-if="isSelect" #input="{ submit, deactivate }">
              <n-select
                v-model:value="selectValue"
                style="width: 200px"
                :placeholder="entity"
                :options="options"
                @update:value="submit($event); selectValue = ''"
                @blur="deactivate"
              />
            </template>
          </n-dynamic-tags>
        </n-space>

        <n-space justify="end" align="center" style="margin-top: 2vh">
          <n-button round secondary type="error" @click="show = false">
            Cancel
          </n-button>
          <n-button round type="success" @click="handleSubmit">
            Submit
          </n-button>
        </n-space>
      </n-card>
    </n-spin>
  </n-modal>
</template>

<script setup>
import { ref, watch, defineAsyncComponent } from 'vue'
import { useMessage, NButton, NIcon, NCard, NModal, NSpace, NDynamicTags, NSpin, NSelect } from 'naive-ui'

const AddOutline = defineAsyncComponent(
  () => import('@vicons/ionicons5/AddOutline')
)
const CloseOutline = defineAsyncComponent(
  () => import('@vicons/ionicons5/CloseOutline')
)

const props = defineProps({
  device: { type: Object, required: true },
  entity: { type: String, default: 'Tags' },
  isSelect: { type: Boolean, default: false },
  options: { type: Array, default: [] },
  value: { type: Array, default: [] }
})
const emits = defineEmits(['save'])
const message = useMessage()

const tags = ref([])
const show = ref(false)
const loading = ref(false)
const selectValue = ref('')

watch(show, (value) => {
  tags.value = (value) ? props.value : []
  loading.value = false
})

function handleSubmit() {
  loading.value = true
  emits(
    'save',
    tags.value,
    () => { show.value = false },
    (msg) => {
      message.error(msg)
      loading.value = false
    }
  )
}
</script>