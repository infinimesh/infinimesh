<template>
  <n-select
    style="max-width: 100%; width: 100%"
    :value="internalValue"
    @update:value="updateInternalValue"
    tag
    filterable
    multiple
    placeholder="Filter devices eg. :uuid:abc"
    :options="filterDeviceOptions"
    :show-arrow="false"
    :render-tag="renderTag"
    class="filter-input"
  />
</template>

<script setup>
import { ref, defineProps, defineEmits, h } from "vue";
import { NSelect, NTag } from "naive-ui";

const { filterTerm } = defineProps(["filterTerm"]);
const emit = defineEmits();

const internalValue = ref(filterTerm);

const updateInternalValue = (value) => {
  internalValue.value = value;
  emit("update:value", value);
};

const renderTag = ({ option, handleClose, ...rest }, i) => {
  return h(
    NTag,
    {
      type: getTagColor(option),
      closable: true,
      onMousedown: (e) => {
        e.preventDefault();
      },
      onClose: (e) => {
        e.stopPropagation();
        handleClose();
      },
    },
    { default: () => option.label }
  );
};

const filterDeviceOptions = [
  {
    label: ":uuid:",
    value: ":uuid:",
    disabled: true,
  },
  {
    label: ":enabled:",
    value: ":enabled:",
    disabled: true,
  },
  {
    label: ":tag:",
    value: ":tag:",
    disabled: true,
  },
  {
    label: ":title:",
    value: ":title:",
    disabled: true,
  },
  {
    label: ":namespace:",
    value: ":namespace:",
    disabled: true,
  },
];

const getTagColor = (option) => {
  const tag = option.label.split(":")[1];

  switch (tag) {
    case "uuid":
      return "primary";
    case "enabled":
      return "success";
    case "tag":
      return "info";
    case "title":
      return "warning";
    case "namespace":
      return "error";
    default:
      return "default";
  }
};
</script>