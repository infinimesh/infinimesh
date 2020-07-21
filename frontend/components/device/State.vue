<template>
  <div id="state">
    <a-row>
      <a-col :span="10">
        <h1 class="lead">{{ title }}</h1>
      </a-col>
      <template v-if="editable">
        <template v-if="active_edit">
          <a-col :span="6">
            <a-button type="primary" icon="close" @click="active_edit = false"
              >Cancel</a-button
            >
          </a-col>
          <a-col :span="6" :offset="1">
            <a-popconfirm
              title="Are you sure saving this state?"
              ok-text="Yes"
              cancel-text="No"
              @confirm="
                handleSaveState();
                active_edit = false;
              "
            >
              <a-button type="success" icon="save">Save</a-button>
            </a-popconfirm>
          </a-col>
        </template>
        <a-col :span="4" :offset="8" v-else>
          <a-button type="primary" icon="edit" @click="active_edit = true"
            >Edit</a-button
          >
        </a-col>
      </template>
    </a-row>
    <p>
      <strong>Version:</strong>
      <u>{{ state.version }}</u>
    </p>
    <p>
      <strong>Updated:</strong>
      {{ ts2str(state.timestamp) }}
    </p>
    <a-spin :spinning="state_updating" v-if="active_edit">
      <a-textarea v-model="desired_state" :autoSize="{ minRows: 3 }" />
    </a-spin>
    <pre v-html="state.data" v-else />
  </div>
</template>

<script>
import Vue from "vue";

export default Vue.component("device-state", {
  props: {
    state: {
      required: true
    },
    title: {
      required: true,
      type: String
    },
    editable: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      today: new Date(),
      active_edit: false,
      desired_state: JSON.stringify(this.state.data, null, 2),
      state_updating: false
    };
  },
  methods: {
    ts2str(ts) {
      let date = ts.match(
        /(?<year>[0-9]{4})-(?<month>[0-9]{2})-(?<day>[0-9]{2})T(?<hour>[0-9]{2}):(?<minutes>[0-9]{2}):(?<seconds>[0-9]{2})(?<miliseconds>\.[0-9]{1,})?Z/
      ).groups;
      let result = "";

      // Day month section
      if (
        Number(date.month) == this.today.getMonth() + 1 &&
        Number(date.year) == this.today.getFullYear()
      ) {
        if (Number(date.day) == this.today.getDate()) {
          result += "Today";
        } else if (Number(date.day) == this.today.getDate() - 1) {
          result += "Yesterday";
        }
      } else {
        result += date.day + "." + date.month;
      }

      // Year section
      if (Number(date.year) != this.today.getFullYear()) {
        result += `.${date.year}`;
      }

      // Time section
      result += ` ${date.hour}:${date.minutes}:${date.seconds}`;

      return result;
    },
    handleSaveState() {
      console.log(this.desired_state);
      this.state_updating = true;
      this.$emit("update", this.desired_state, () => {
        this.state_updating = false;
      });
    }
  }
});
</script>

<style lang="less" scoped>
pre {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 10px;
  width: 95%;
}
textarea {
  background: @infinimesh-dark-purple;
  width: 95%;
  font-family: monospace, monospace;
}
button {
  margin-top: 0.3rem;
}
</style>
