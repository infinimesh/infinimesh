<template>
  <nuxt-link :to="{ name: 'dashboard-devices-id', params: { id: device.id } }">
    <a-dropdown :trigger="['contextmenu']">
      <a-card
        :hoverable="true"
        :bordered="false"
        :ref="`device-card-${device.id}`"
        :class="selected ? 'card-selected' : ''"
      >
        <template slot="title">{{ device.name }}</template>
        <template slot="extra">
          <b class="muted">{{ device.id }}</b>
          <a-tooltip
            :title="device.enabled ? 'Device enabled' : 'Device is not enabled'"
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
          <a-row v-else type="flex" justify="center" class="muted"
            >No tags were provided</a-row
          >
        </template>
      </a-card>

      <a-menu slot="overlay">
        <a-menu-item key="open">
          <nuxt-link
            :to="{ name: 'dashboard-devices-id', params: { id: device.id } }"
          >
            <a-button type="link"> Open </a-button>
          </nuxt-link>
        </a-menu-item>
        <a-menu-item key="toogle">
          <a-button type="link" @click="handleToogleDevice(false)">
            {{ device.enabled ? "Disable" : "Enable" }}
          </a-button>
        </a-menu-item>
        <a-menu-item key="toogle-selection">
          <a-button
            type="link"
            @click="$emit((selected ? 'de' : '') + 'select', device.id)"
          >
            {{ selected ? "Deselect" : "Select" }}
          </a-button>
        </a-menu-item>
      </a-menu>
    </a-dropdown>
  </nuxt-link>
</template>

<script>
import Vue from "vue";
import deviceControlMixin from "@/mixins/device-control";

export default Vue.component("device-list-card", {
  mixins: [deviceControlMixin],
  props: {
    device: {
      required: true,
      type: Object,
    },
    selected: {
      required: false,
      default: false,
      type: Boolean,
    },
  },
});
</script>

<style scoped>
.card-selected {
  -webkit-box-shadow: 20px 15px 10px 5px rgba(0, 0, 0, 0.7);
  -moz-box-shadow: 20px 15px 10px 5px rgba(0, 0, 0, 0.7);
  box-shadow: 20px 15px 10px 5px rgba(0, 0, 0, 0.7);
}
</style>