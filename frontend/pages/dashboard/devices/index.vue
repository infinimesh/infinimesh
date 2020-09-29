<template>
  <div id="devicesTable">
    <a-row :gutter="{ md: 10, lg: 10, xl: 10, xxl: 10 }" type="flex" id="root">
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
          <a-row
            class="create-form"
            v-if="device.type && device.type == 'create-form'"
            type="flex"
            justify="center"
            align="middle"
          >
            <a-icon
              type="plus"
              style="font-size: 6rem; height: 100%; width: 100%"
              @click="addDeviceActive = true"
            />
            <device-add
              :active="addDeviceActive"
              @cancel="addDeviceActive = false"
              @add="handleDeviceAdd"
            />
          </a-row>
          <device-list-card :device="device" v-else />
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import DeviceAdd from "@/components/device/Add.vue";
import DeviceListCard from "@/components/device/ListCard.vue";

export default {
  name: "devicesTable",
  components: {
    DeviceAdd,
    DeviceListCard,
  },
  data() {
    return {
      addDeviceActive: false,
    };
  },
  computed: {
    pool: {
      deep: true,
      get() {
        return this.$store.state.devices.pool;
      },
    },
    poolCols: {
      deep: true,
      get() {
        let pool = [{ type: "create-form" }, ...this.pool];
        if (!pool.length) return pool;
        let div = 1;
        switch (this.gridSize) {
          case "xs": {
            div = 1;
            break;
          }
          case "sm": {
            div = 2;
            break;
          }
          case "md": {
            div = 2;
            break;
          }
          case "lg": {
            div = 3;
            break;
          }
          case "xl": {
            div = 3;
            break;
          }
          case "xxl": {
            div = 4;
            break;
          }
        }
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
    gridSize() {
      return this.$store.state.window.gridSize;
    },
  },
  methods: {
    handleDeviceAdd(device) {
      this.$store.dispatch("devices/add", {
        device: device,
        error: (err) => {
          this.$notification.error({
            message: "Failed to create the device",
            description: `Response: ${err.response.data.message}`,
            placement: "bottomRight",
            duration: 10,
          });
        },
        always: () => {
          this.addDeviceActive = false;
        },
      });
    },
  },
};
</script>

<style scoped>
.create-form {
  border-radius: var(--border-radius-base);
  background: var(--primary-color)-dark;
  border: var(--border-base);
  min-height: 8rem;
  cursor: pointer;
}
.create-form .anticon {
  color: var(--icon-color-dark);
}
</style>