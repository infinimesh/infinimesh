<template>
  <div id="state">
    <a-row>
      <a-col :span="10">
        <h1 class="lead lead-dark">{{ title }}</h1>
      </a-col>
      <a-col :span="4" :offset="8" v-if="editable && !active_edit">
        <a-button type="primary" icon="edit" @click="active_edit = true"
          >Edit</a-button
        >
      </a-col>
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
    <a-row v-if="active_edit">
      <a-col :span="10">
        <a-button type="primary" icon="close" @click="active_edit = false"
          >Cancel</a-button
        >
      </a-col>
      <a-col :span="10" :offset="2" style="text-align: right">
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
    </a-row>
  </div>
</template>

<script>
import Vue from "vue";

const formatDateNumber = (num, n = 2) => {
  num = num.toString();
  while (num.length < n) {
    num = "0" + num;
  }
  return num;
};
const date2Object = date => {
  return {
    day: date.getDate(),
    month: date.getMonth(),
    year: date.getFullYear(),
    hour: date.getHours(),
    minute: date.getMinutes(),
    second: date.getSeconds()
  };
};

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
      active_edit: false,
      desired_state: JSON.stringify(this.state.data, null, 2),
      state_updating: false
    };
  },
  methods: {
    ts2str(ts) {
      let today = date2Object(new Date());

      let date = new Date(Date.parse(ts));
      date = date2Object(date);

      let result = "";

      // Day month section
      if (
        Number(date.month) == today.month &&
        Number(date.year) == today.year
      ) {
        if (Number(date.day) == today.day) {
          result += "Today";
        } else if (Number(date.day) == today.day - 1) {
          result += "Yesterday";
        }
      } else {
        result +=
          formatDateNumber(date.day) + "." + formatDateNumber(date.month + 1);
      }

      // Year section
      if (Number(date.year) != today.year) {
        result += `.${formatDateNumber(date.year, 4)}`;
      }

      // Time section
      result += ` ${formatDateNumber(date.hour)}:${formatDateNumber(
        date.minute
      )}:${formatDateNumber(date.second)}`;

      return result;
    },
    handleSaveState() {
      this.state_updating = true;
      this.$emit("update", this.desired_state, () => {
        this.state_updating = false;
      });
    }
  }
});
</script>

<style scoped>
pre {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 10px;
  width: 95%;
}
textarea {
  background: var(--primary-color)-dark;
  width: 95%;
  font-family: monospace, monospace;
}
button {
  margin-top: 0.3rem;
}
strong {
  color: #d1d1ff;
}
</style>
