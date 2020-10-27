<template>
  <div id="devicesTable">
    <a-row
      type="flex"
      style="padding: 4px 10px"
      class="tile-bar tile-bar-transparent"
      align="middle"
    >
      <template>
        <a-col>
          <a-input
            placeholder="Search device..."
            class="devices-search-input"
          />
        </a-col>
        <a-col>
          <a-row type="flex" justify="center">
            <a-col>
              <a-switch
                id="group-by-tags-switch"
                un-checked-children="Group by tags"
                checked-children="Whole registry"
                v-model="groupByTags"
              />
            </a-col>
          </a-row>
        </a-col>
      </template>

      <template v-if="selectedDevices.length">
        <a-col>
          <a-button type="link" @click="selectedDevices = []" icon="close"
            >Deselect all
          </a-button>
        </a-col>
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
      </template>
      <template v-else>
        <a-col>
          <a-button
            type="success"
            @click="selectedDevices = pool.map((d) => d.id)"
            >Select All
          </a-button>
        </a-col>
      </template>
    </a-row>
    <template v-if="groupByTags">
      <a-collapse :bordered="false" accordion class="tags-collapse-tiles-wrap">
        <a-collapse-panel
          v-for="tag in tags"
          :key="tag"
          class="tags-collapse-tile-wrap"
        >
          <span
            slot="header"
            style="color: var(--text-color); font-size: 18px"
            >{{ tag }}</span
          >
          <div slot="extra" @click="(e) => e.stopPropagation()">
            <a-button
              type="success"
              style="border-radius: 100px; height: 24px"
              @click="
                selectedDevices.push(
                  ...pool.filter((d) => d.tags.includes(tag)).map((d) => d.id)
                )
              "
              >Select All
            </a-button>
          </div>

          <device-pool
            :div="div"
            :selected="selectedDevices"
            :pool="pool.filter((d) => d.tags.includes(tag))"
            :grouped="true"
            @select="(id) => selectedDevices.push(id)"
            @deselect="
              (id) => selectedDevices.splice(selectedDevices.indexOf(id), 1)
            "
            @select-all="selectedDevices = pool.map((d) => d.id)"
            style="
              background-color: var(--secondary-color);
              border-radius: var(--border-radius-base);
              padding-bottom: 10px;
            "
          />
        </a-collapse-panel>
      </a-collapse>
    </template>
    <device-pool
      :div="div"
      :selected="selectedDevices"
      :pool="mainPool"
      @select="(id) => selectedDevices.push(id)"
      @deselect="(id) => selectedDevices.splice(selectedDevices.indexOf(id), 1)"
      @select-all="selectedDevices = pool.map((d) => d.id)"
    >
      <a-row
        class="create-form"
        type="flex"
        justify="center"
        align="middle"
        slot="device-create-form"
      >
        <nuxt-link
          v-if="user.default_namespace.id === namespace"
          :to="{ name: 'dashboard-namespaces', query: { create: true } }"
          no-prefetch
        >
          <h3 style="padding: 15px">
            <p>
              You can't create devices in your root namespace, switch to another
              one to perform device create.
            </p>
            <p>
              Click here to create new namespace, or switch namespace on top of
              the page.
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
    </device-pool>
  </div>
</template>

<script>
import DevicePool from "@/components/device/Pool.vue";
import DeviceAdd from "@/components/device/Add.vue";

const divs = {
  xs: 1,
  sm: 2,
  md: 2,
  lg: 3,
  xl: 3,
  xxl: 4,
};

export default {
  name: "devicesTable",
  components: {
    DeviceAdd,
    DevicePool,
  },
  data() {
    return {
      addDeviceActive: false,
      selectedDevices: [],
      groupByTags: false,
    };
  },
  watch: {
    selectedDevices() {
      this.selectedDevices.filter((e, i, self) => self.indexOf(e) === i);
    },
  },
  computed: {
    tags() {
      return this.pool.reduce((r, el) => {
        el.tags.forEach((t) => r.add(t));
        return r;
      }, new Set());
    },
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
    mainPool: {
      deep: true,
      get() {
        return [{ type: "create-form" }, ...this.pool];
      },
    },
    div() {
      return divs[this.$store.state.window.gridSize];
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
  min-height: 32px;
}
.tile-bar-transparent {
  background: none;
}
.tile-bar .ant-btn-link {
  color: var(--line-color) !important;
}
.tile-bar .ant-btn {
  height: 90%;
  border-radius: 100px;
}
#group-by-tags-switch {
  background-color: var(--switch-color);
}
#group-by-tags-switch.ant-switch-checked {
  background-color: var(--success-color);
}
.devices-search-input {
  height: 90%;
  border: 1px solid var(--primary-color);
  border-radius: 100px;
  min-width: 256px;
}
@media (max-width: 768px) {
  .tile-bar > [class*="ant-col"] {
    margin-top: 3px;
    margin-bottom: 3px;
  }
}
.tile-bar > .ant-col + .ant-col {
  margin-left: 15px;
}
.tags-collapse-tiles-wrap {
  margin-top: 10px;
  margin-bottom: 5px;
}
.tags-collapse-tile-wrap {
  border-radius: var(--border-radius-base);
  background: var(--primary-color);
}
.tags-collapse-tile-wrap:last-child {
  border-radius: var(--border-radius-base);
}
.ant-collapse-borderless {
  background-color: var(--secondary-color) !important;
}
</style>