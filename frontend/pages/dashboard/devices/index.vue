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
          <nuxt-link
            :to="{ name: 'dashboard-devices-id', params: { id: device.id } }"
          >
            <a-card :hoverable="true" :bordered="false">
              <template slot="title">
                {{ device.name }}
              </template>
              <template slot="extra">
                <b class="muted">
                  {{ device.id }}
                </b>
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
                <a-row v-else type="flex" justify="center" class="muted">
                  No tags were provided
                </a-row>
              </template>
            </a-card>
          </nuxt-link>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script>
export default {
  name: "devicesTable",
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
        console.log("DIV:", div);
        if (div == 1) {
          return [this.pool];
        }
        let res = new Array(div);
        for (let i = 0; i < div; i++) {
          res[i] = new Array();
        }
        for (let i = 0; i <= this.pool.length; i++) {
          for (let j = 0; j < div && i + j < this.pool.length; j++) {
            res[j].push(this.pool[i + j]);
          }
          i += div - 1;
        }
        return res;
      }
    },
    gridSize() {
      return this.$store.state.window.gridSize;
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
</style>
