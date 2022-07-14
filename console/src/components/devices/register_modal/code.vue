<template>
  <n-space justify="center" align="center">
    <n-input v-for="n in 3" minlength="1" maxlength="1" placeholder="" :ref="el => chars_refs[n - 1] = el"
      @update:value="v => update(n - 1, v)" :value="chars[n - 1]" class="char-input" size="large" />
    <span>-</span>
    <n-input v-for="n in 3" minlength="1" maxlength="1" placeholder="" :ref="el => chars_refs[n + 2] = el"
      @update:value="v => update(n + 2, v)" :value="chars[n + 2]" class="char-input" size="large" />
  </n-space>
</template>

<script setup>
import { ref, watch, onBeforeUpdate } from "vue"
import { NSpace, NInput } from "naive-ui"

const chars_refs = ref([])

const chars = ref(new Array(6))

function update(id, value) {
  console.log(id, value)
  chars.value[id] = value.toUpperCase()
  try {
    let input;
    if (value == "" && id > 0) { // char removed after the first input
      input = chars_refs.value[id - 1] // focus on the previous one
    } else if (value != "" && id < 5) { // char enterd anywhere except the last input
      input = chars_refs.value[id + 1] // focus on the next one
    } else {
      return
    }
    input.focus()
  } catch (e) {
    console.error(e)
    console.log("this might help", id, chars_refs)
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

</script>

<style>
.char-input {
  text-align: center;
  max-width: calc(3 * var(--n-font-size));
}
</style>