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
            :style="deviceCreateFormStyle"
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
          <nuxt-link v-else :to="{ name: 'dashboard-devices-id', params: { id: device.id } }">
            <a-card :hoverable="true" :bordered="false" :ref="`device-card-${device.id}`">
              <template slot="title">{{ device.name }}</template>
              <template slot="extra">
                <b class="muted">{{ device.id }}</b>
                <a-tooltip
                  :title="
                    device.enabled ? 'Device enabled' : 'Device is not enabled'
                  "
                  placement="bottom"
                >
                  <a-icon
                    type="bulb"
                    :style="{ color: device.enabled ? '#52c41a' : '#eb2f96' }"
                    theme="filled"
                  />
                </a-tooltip>
              </template>
              <template>
                <a-row v-if="device.tags.length">
                  Tags:
                  <a-tag v-for="tag in device.tags" :key="tag">{{ tag }}</a-tag>
                </a-row>
                <a-row v-else type="flex" justify="center" class="muted">No tags were provided</a-row>
              </template>
            </a-card>
          </nuxt-link>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import DeviceAdd from "@/components/device/Add.vue";

export default {
  name: "devicesTable",
  components: {
    DeviceAdd
  },
  data() {
    return {
      addDeviceActive: false
    };
  },
  computed: {
    pool: {
      deep: true,
      get() {
        return this.$store.state.devices.pool;
      }
    },
    poolCols: {
      deep: true,
      get() {
        if (!this.pool.length) return this.pool;
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
        let pool = [{ type: "create-form" }, ...this.pool];
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
      }
    },
    deviceCreateFormStyle() {
      return {
        "--device-card-height": this.deviceCardHeight
      };
    },
    deviceCardHeight: {
      deep: true,
      get() {
        if (this.$refs.length) {
          return this.$refs.reduce((curr, el) => {
            if (el < curr) curr = el.clientHeight;
            return curr;
          }, 1000);
        } else {
          return "8rem";
        }
      }
    },
    gridSize() {
      return this.$store.state.window.gridSize;
    }
  },
  methods: {
    handleDeviceAdd(device) {
      console.log(device);
      this.$store.dispatch("devices/add", device);
      this.addDeviceActive = false;
    }
  }
};
</script>

<style scoped>
#root {
  padding: 10px;
}
</style>
<style lang="less" scoped>
.muted {
  color: @infinimesh-dark-purple;
}
.create-form {
  border-radius: @border-radius-base;
  background: @infinimesh-dark-purple;
  border: 1px dashed white;
  min-height: var(--device-card-height);
  cursor: pointer;
}
</style>
