<template>
  <div id="devicesTable">
    <div class="tile-bar">
      <a-row type="flex" align="middle" style="padding: 0 10px">
        <template>
          <a-col :span="5">
            <a-input placeholder="Search device..." style="height: 90%" />
          </a-col>
          <div role="separator" class="tile-bar-vertical-divider"></div>
          <a-col>
            <a-row type="flex" justify="center">
              <a-col>
                <a-switch
                  id="group-by-tags-switch"
                  un-checked-children="Group by tags"
                  checked-children="Whole registry"
                />
              </a-col>
            </a-row>
          </a-col>
          <div role="separator" class="tile-bar-vertical-divider"></div>
        </template>
        <template v-if="selectedDevices.length">
          <a-col>
            <a-button
              type="success"
              style="margin-right: 5px"
              @click="toogleAll(true)"
              >Enable All
            </a-button>
            <a-button type="danger" @click="toogleAll(false)"
              >Disable All
            </a-button>
          </a-col>
          <div role="separator" class="tile-bar-vertical-divider"></div>
          <a-col>
            <a-button type="link" @click="selectedDevices = []" icon="close"
              >Deselect all
            </a-button>
          </a-col>
          <div role="separator" class="tile-bar-vertical-divider"></div>
        </template>
      </a-row>
    </div>
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
            <nuxt-link
              v-if="user.default_namespace.id === namespace"
              :to="{ name: 'dashboard-namespaces', query: { create: true } }"
              no-prefetch
            >
              <h3 style="padding: 15px">
                <p>
                  You can't create devices in your root namespace, switch to
                  another one to perform device create.
                </p>
                <p>
                  Click here to create new namespace, or switch namespace on top
                  of the page.
                </p>
              </h3>
            </nuxt-link>
            <template v-else>
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
            </template>
          </a-row>
          <device-list-card
            v-else
            :device="device"
            :selected="selectedDevices.includes(device.id)"
            @select="(id) => selectedDevices.push(id)"
            @deselect="
              (id) => selectedDevices.splice(selectedDevices.indexOf(id), 1)
            "
            @select-all="selectedDevices = pool.map((d) => d.id)"
          />
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
      selectedDevices: [],
    };
  },
  computed: {
    user() {
      return this.$store.getters.loggedInUser;
    },
    namespace() {
      return this.$store.getters["devices/currentNamespace"];
    },
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
            duration: 10,
          });
        },
        always: () => {
          this.addDeviceActive = false;
        },
      });
    },
    toogleAll(enable) {
      let vm = this;
      this.updateAll(
        () => {
          return {
            enabled: enable,
          };
        },
        {
          success: (response) => {
            let res = response.reduce(
              (r, el) => {
                r[el.status == 200 ? "success" : "fail"]++;
                return r;
              },
              {
                success: 0,
                fail: 0,
              }
            );
            vm.$notification.info({
              message: `${enable ? "Enabled" : "Disabled"} ${
                response.length
              } devices`,
              description: `Result: success - ${res.success}, failed: ${res.fail}.`,
            });
          },
        }
      );
    },
    updateAll(modifier, { success, error }) {
      let patchPromises = this.pool
        .filter((d) => this.selectedDevices.includes(d.id))
        .map((device) => {
          return this.$axios({
            url: `/api/devices/${device.id}`,
            method: "patch",
            data: modifier(device),
          });
        });
      Promise.all(patchPromises)
        .then((res) => {
          if (success) success(res);
        })
        .catch((err) => {
          if (error) error(err);
        })
        .then(() => {
          this.$store.dispatch("devices/get");
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
.tile-bar {
  margin-top: 10px;
  width: 100%;
  background: var(--primary-color);
  border-radius: var(--border-radius-base);
  color: var(--line-color);
}
.tile-bar .ant-btn-link {
  color: white !important;
}
.tile-bar .ant-btn {
  height: 90%;
}
.tile-bar-vertical-divider {
  box-sizing: border-box;
  margin: 0 10px;
  padding: 0;
  color: rgba(0, 0, 0, 0.65);
  font-size: 14px;
  font-variant: tabular-nums;
  line-height: 1.5;
  list-style: none;
  font-feature-settings: "tnum", "tnum";
  background: #e8e8e8;
  vertical-align: middle;
  width: 1px;
  min-height: 32px;
}
#group-by-tags-switch {
  /* background-color: rgb(160, 160, 160); */
  background-color: var(--switch-color);
}
#group-by-tags-switch.ant-switch-checked {
  background-color: var(--success-color);
}
</style>