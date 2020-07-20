<template>
  <div id="device">
    <a-spin
      :spinning="!device || !device.name"
      style="min-height: 15rem"
      size="large"
      tip="Loading device..."
    >
      <a-row style="padding-top: 10px">
        <a-col :xxl="{ span: 10, offset: 1 }">
          <transition name="fade">
            <h1 class="lead" v-if="device">
              {{ device.name }} <span class="muted">{{ device.id }}</span>
            </h1>
          </transition>
        </a-col>
        <a-col :xxl="{ span: 1, offset: 1 }">
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
        <a-col :xxl="{ span: 12, offset: 1 }">
          <transition-group name="slide">
            <a-card title="Details" key="details" v-if="device">
              <template>
                <a-row v-if="device.tags && device.tags.length">
                  <p>
                    Tags:
                    <a-tag v-for="tag in device.tags" :key="tag">{{
                      tag
                    }}</a-tag>
                  </p>
                </a-row>
                <a-row v-else type="flex" justify="center" class="muted">
                  <p>No tags were provided</p>
                </a-row>
              </template>
              <template>
                <p>
                  Namespace: <u>{{ device.namespace }}</u>
                </p>
              </template>
            </a-card>
            <a-card title="State" key="state" v-if="device && device.state">
              <a-row>
                <a-col :span="12">
                  <device-state
                    title="Reported"
                    :state="device.state.shadow.reported"
                  />
                </a-col>
                <a-col :span="12">
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
import DeviceState from "@/components/device/DeviceState";

export default {
  components: { DeviceState },
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
      id: this.$route.params.id
    };
    this.$axios.get(`/devices/${this.device.id}`).then(res => {
      this.device = res.data.device;
    });
    this.deviceStateGet();
  },
  methods: {
    async deviceStateGet() {
      await this.$axios.get(`/devices/${this.device.id}/state`).then(res => {
        this.device = {
          ...this.device,
          state: res.data
        };
      });
    },
    handleStateUpdate(state, callback) {
      this.$axios({
        url: `/devices/${this.device.id}/state`,
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
.lead {
  font-size: 2rem;
  color: #fff;
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
</style>
