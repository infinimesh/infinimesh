<template>
  <div id="device">
    <a-spin
      :spinning="!device || !device.name"
      style="min-height: 15rem"
      size="large"
      tip="Loading device..."
    >
      <a-row style="padding-top: 10px">
        <a-col
          :xs="{ span: 21 }"
          :sm="{ span: 20, offset: 1 }"
          :md="{ span: 20, offset: 1 }"
          :lg="{ span: 20, offset: 1 }"
          :xl="{ span: 16, offset: 1 }"
          :xxl="{ span: 12, offset: 1 }"
        >
          <transition name="fade">
            <h1 class="lead" v-if="device">
              {{ device.name }}
              <span class="muted">{{ device.id }}</span>
            </h1>
          </transition>
        </a-col>
        <a-col
          :xs="1"
          :sm="1"
          :md="1"
          :xl="{ span: 1, offset: 1 }"
          :xxl="{ span: 1, offset: 1 }"
        >
          <a-row type="flex" justify="end">
            <a-tooltip
              :title="
                device.enabled ? 'Device enabled' : 'Device is not enabled'
              "
              placement="right"
            >
              <transition name="fade">
                <a-icon
                  v-show="device"
                  type="bulb"
                  class="device-state-bulb"
                  :style="{ color: deviceStateBulbColor }"
                  theme="filled"
                />
              </transition>
            </a-tooltip>
          </a-row>
        </a-col>
      </a-row>
      <a-row>
        <a-col
          :sm="{ span: 22, offset: 1 }"
          :md="{ span: 22, offset: 1 }"
          :lg="{ span: 22, offset: 1 }"
          :xl="{ span: 18, offset: 1 }"
          :xxl="{ span: 14, offset: 1 }"
        >
          <transition-group name="slide">
            <a-card title="Details" key="details" v-if="device" hoverable>
              <template>
                <a-row v-if="device.tags && device.tags.length">
                  <p>
                    Tags:
                    <a-tag v-for="tag in device.tags" :key="tag">
                      {{ tag }}
                    </a-tag>
                  </p>
                </a-row>
                <a-row v-else type="flex" justify="center" class="muted">
                  <p>No tags were provided</p>
                </a-row>
              </template>
              <template>
                <p>
                  Namespace:
                  <u>{{ device.namespace }}</u>
                </p>
              </template>
            </a-card>
            <a-card title="Actions" key="actions" v-if="device" hoverable>
              <device-actions
                :device-id="device.id"
                @delete="handleDeviceDelete"
              />
            </a-card>
            <a-card
              title="State"
              key="state"
              v-if="device && device.state"
              hoverable
            >
              <a-row>
                <a-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12" :xxl="12">
                  <device-state
                    title="Reported"
                    :state="device.state.shadow.reported"
                  />
                </a-col>
                <a-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12" :xxl="12">
                  <device-state
                    title="Desired"
                    :state="device.state.shadow.desired"
                    :editable="true"
                    @update="handleStateUpdate"
                  />
                </a-col>
              </a-row>
            </a-card>
          </transition-group>
        </a-col>
      </a-row>
    </a-spin>
  </div>
</template>

<script>
import DeviceState from "@/components/device/State";
import DeviceActions from "@/components/device/Actions";

export default {
  /**
   * Represents Device as both object and component
   * @displayName Device
   */
  components: { DeviceState, DeviceActions },
  props: {
    /**
     * Device ID - not required if component is mounted via Router _id
     */
    deviceId: {
      required: false
    }
  },
  data() {
    return {
      deviceObject: false
    };
  },
  computed: {
    device: {
      get() {
        return this.deviceObject;
      },
      set(obj) {
        this.deviceObject = { ...this.deviceObject, ...obj };
      }
    },
    deviceStateBulbColor() {
      if (!(this.device && this.device.enabled !== undefined)) {
        return "black";
      } else if (this.device.enabled) {
        return "#52c41a";
      } else {
        return "#eb2f96";
      }
    }
  },
  mounted() {
    this.device = {
      id: this.deviceId || this.$route.params.id
    };
    // Getting Device data from API
    this.$axios
      .get(`/api/devices/${this.device.id}`)
      .then(res => {
        this.device = res.data.device;
        this.socket = new WebSocket(
          `wss://${this.$config.baseURL.replace("https://", "")}/devices/${
            this.device.id
          }/state/stream`,
          ["Bearer", this.$auth.getToken("local").split(" ")[1]]
        );
        this.socket.onmessage = msg => {
          this.device.state.shadow.reported = JSON.parse(
            msg.data
          ).result.reportedState;
        };
        window.addEventListener("beforeunload", function(event) {
          this.socket.close();
        });
      })
      .catch(res => {
        if (res.response.status == 404) {
          this.$notification.error({
            message: "Device wasn't found",
            description: "Redirecting...",
            placement: "bottomRight"
          });
        } else if (res.response.status == 403) {
          this.$notification.error({
            message: "You have no access to this device",
            description: "Redirecting...",
            placement: "bottomRight"
          });
        }
        this.$router.push({ name: "dashboard-devices" });
      });
    this.deviceStateGet();
  },
  methods: {
    /**
     * Obtains device state(s) (desired and reported) and merges them into deviceObject
     */
    async deviceStateGet() {
      await this.$axios
        .get(`/api/devices/${this.device.id}/state`)
        .then(res => {
          this.device = {
            ...this.device,
            state: res.data
          };
        });
    },
    /**
     * Performs PATCH /device/id and changes desired state to given.
     * Used by device-state component so it invokes callback after patch response handling is done
     * @param {String} state, the desired state as in JSON String
     */
    handleStateUpdate(state, callback) {
      this.$axios({
        url: `/api/devices/${this.device.id}/state`,
        method: "patch",
        data: state
      })
        .then(res => {
          this.deviceStateGet();
        })
        .catch(res => {
          console.log(res);
        })
        .then(() => {
          callback();
        });
    },
    handleDeviceDelete() {
      this.$axios({
        url: `/api/devices/${this.device.id}`,
        method: "delete"
      }).then(() => {
        this.$message.success("Device successfuly deleted!");
        this.$store.dispatch("devices/get");
        this.$router.push({ name: "dashboard-devices" });
      });
    }
  },
  validate({ params }) {
    return /0[xX][0-9a-fA-F]+/.test(params.id);
  }
};
</script>

<style>
#device {
  overflow: hidden;
  font-family: Exo;
  font-weight: 500;
}
.muted {
  opacity: 0.7;
}

.slide-leave-active,
.slide-enter-active {
  transition: 1s;
}
.slide-enter {
  transform: translate(100%, 0);
}
.slide-leave-to {
  transform: translate(-100%, 0);
}
.device-state-bulb {
  font-size: 1.5rem;
  padding-top: 0.8rem;
}

.ant-card + .ant-card {
  margin-top: 1rem;
}
</style>

<style lang="less" scoped>
.lead {
  font-size: 2rem;
  color: @line-color;
}
</style>
