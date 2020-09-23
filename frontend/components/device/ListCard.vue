<template>
  <nuxt-link :to="{ name: 'dashboard-devices-id', params: { id: device.id } }">
    <a-card
      :hoverable="true"
      :bordered="false"
      :ref="`device-card-${device.id}`"
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
  </nuxt-link>
</template>

<script>
import Vue from "vue";

export default Vue.component("device-list-card", {
  props: {
    device: {
      required: true,
      type: Object
    }
  }
});
</script>
