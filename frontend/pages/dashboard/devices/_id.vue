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
            <h1 class="lead" v-if="device && !active_edit">
              {{ device.name }}
              <span class="muted">{{ device.id }}</span>
            </h1>
            <a-input
              v-else-if="device && active_edit"
              placeholder="Enter new device name"
              class="device-name-input"
              v-model="device.name"
            />
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
            <a-card key="details" v-if="device" hoverable>
              <a-row slot="title" type="flex" justify="space-between">
                <a-col :span="3"> Details </a-col>
                <a-col :xxl="5" :lg="6" :md="9" :sm="10" v-if="active_edit">
                  <a-space>
                    <a-button
                      type="primary"
                      icon="close"
                      @click="active_edit = false"
                      >Cancel</a-button
                    >
                    <a-button type="success" icon="save" @click="patchDevice"
                      >Save</a-button
                    >
                  </a-space>
                </a-col>
                <a-col :lg="2" :md="3" :sm="4" v-else>
                  <a-button
                    type="primary"
                    icon="edit"
                    @click="active_edit = true"
                    >Edit</a-button
                  >
                </a-col>
              </a-row>
              <template v-if="active_edit">
                <a-select
                  mode="tags"
                  :token-separators="[',']"
                  v-model="device.tags"
                  style="min-width: 50%; margin: 15px 0"
                  placeholder="Enter a comma-separated list of tags, e.g. tag1, tag2"
                />
              </template>
              <template v-else>
                <a-row v-if="device.tags && device.tags.length">
                  <a-col :span="2">Tags:</a-col>
                  <a-col :span="22">
                    <a-tag
                      v-for="tag in device.tags"
                      :key="tag"
                      style="margin-bottom: 5px"
                      >{{ tag }}</a-tag
                    >
                  </a-col>
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
                :device="device"
                @delete="handleDeviceDelete"
                @toogle="handleToogleDevice"
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

import deviceControlMixin from "@/mixins/device-control";

export default {
  /**
   * Represents Device as both object and component
   * @displayName Device
   */
  components: { DeviceState, DeviceActions },
  mixins: [deviceControlMixin],
  props: {
    /**
     * Device ID - not required if component is mounted via Router _id
     */
    deviceId: {
      required: false,
    },
  },
  data() {
    return {
      deviceObject: false,
      active_edit: false,
    };
  },
  computed: {
    device: {
      get() {
        return this.deviceObject;
      },
      set(obj) {
        this.deviceObject = { ...this.deviceObject, ...obj };
      },
    },
    deviceStateBulbColor() {
      if (!(this.device && this.device.enabled !== undefined)) {
        return "black";
      } else if (this.device.enabled) {
        return "#52c41a";
      } else {
        return "#eb2f96";
      }
    },
  },
  mounted() {
    this.device = {
      id: this.deviceId || this.$route.params.id,
    };
    // Getting Device data from API
    this.refresh();
    this.$store.commit("window/setTopAction", {
      icon: "undo",
      callback: this.refresh,
    });
  },
  beforeDestroy() {
    this.$store.commit("window/unsetTopAction");
  },
  methods: {
    patchDevice() {
      this.active_edit = false;
      this.handleDeviceUpdate(
        {
          name: this.device.name,
          tags: this.device.tags,
        },
        {
          refresh: true,
          success: () => {
            this.$message.success(`Device successfuly updated!`);
          },
          error: (e) => {
            this.$notification.error({
              message: `Error updating device`,
              description: e.response.data.message,
            });
          },
        }
      );
    },
    validate({ params }) {
      return /0[xX][0-9a-fA-F]+/.test(params.id);
    },
  },
};
</script>

<style scoped>
.device-name-input {
  width: 50%;
  margin-bottom: 16px;
  font-size: 1.7rem;
  padding: 0 !important;
  line-height: normal !important;
}
</style>
<style>
#device {
  overflow: hidden;
  font-family: Exo;
  font-weight: 500;
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
