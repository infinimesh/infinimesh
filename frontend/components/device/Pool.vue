<template>
  <a-row :gutter="{ md: 10, lg: 10, xl: 10, xxl: 10 }" type="flex">
    <a-col
      :xs="{ span: 24 }"
      :ms="{ span: 12 }"
      :md="{ span: 12 }"
      :lg="{ span: 8 }"
      :xl="{ span: 8 }"
      :xxl="{ span: 6 }"
      v-for="(col, i) in poolCols"
      :key="i"
    >
      <div style="padding-top: 10px" v-for="device in col" :key="device.id">
        <slot
          name="device-create-form"
          v-if="device.type && device.type == 'create-form'"
        />
        <device-list-card
          v-else
          :device="device"
          :selected="selected.includes(device.id)"
          @select="(id) => $emit('select', id)"
          @deselect="(id) => $emit('deselect', id)"
          @select-all="$emit('select-all')"
          @tag-clicked="(tag) => $emit('tag-clicked', tag)"
        />
      </div>
    </a-col>
  </a-row>
</template>

<script>
import DeviceListCard from "@/components/device/ListCard.vue";

export default {
  name: "device-pool",
  components: {
    DeviceListCard,
  },
  props: {
    div: {
      required: true,
      type: Number,
    },
    selected: {
      required: true,
      type: Array,
    },
    pool: {
      required: true,
      type: Array,
    },
  },
  computed: {
    poolCols: {
      deep: true,
      get() {
        return this.splitIntoCols(this.div, this.pool);
      },
    },
  },
  methods: {
    splitIntoCols(div, pool) {
      if (!pool.length) return pool;
      if (div == 1) {
        return [pool];
      }
      let res = new Array(div);
      for (let i = 0; i < div; i++) {
        res[i] = new Array();
      }
      for (let i = 0; i <= pool.length; i++) {
        for (let j = 0; j < div && i + j < pool.length; j++) {
          res[j].push(pool[i + j]);
        }
        i += div - 1;
      }
      return res;
    },
  },
};
</script>
