<template>
  <n-space justify="center" align="center">
    <n-input v-for="n in 3" minlength="1" maxlength="1" style="max-width: 3.5vw" placeholder="" :ref="'char' + (n - 1)"
      @update:value="v => update(n - 1, v)" :value="chars[n - 1]" />
    <span>-</span>
    <n-input v-for="n in 3" minlength="1" maxlength="1" style="max-width: 3.5vw" placeholder="" :ref="'char' + (2 + n)"
      @update:value="v => update(n + 2, v)" :value="chars[n + 2]" />
  </n-space>
</template>

<script setup>
import { ref, watch } from "vue"
import { NSpace, NInput } from "naive-ui"

const char0 = ref(null)
const char1 = ref(null)
const char2 = ref(null)
const char3 = ref(null)
const char4 = ref(null)
const char5 = ref(null)

const chars_refs = ref([char0, char1, char2, char3, char4, char5])
const chars = ref(new Array(6))

function update(id, value) {
  chars.value[id] = value.toUpperCase()

  if (value == "" && id > 0) {
    chars_refs.value[id - 1].value[0].focus()
  } else if (value != "" && id < 5) {
    chars_refs.value[id + 1].value[0].focus()
  }
}

const emit = defineEmits(['update:nextEnabled', 'update:value'])

watch(chars, () => {
  let r = chars.value.join('')
  if (r && r.length == 6) {
    emit('update:nextEnabled', true)
    emit('update:value', {
      code: r,
    })
  } else {
    emit('update:nextEnabled', false)
  }
}, {
  deep: true
})

window.addEventListener("keydown", function (e) {
  if (e.keyCode == 8) {
    console.log('go back')
  }
});

</script>