<template>
  <a-modal
    :visible="active"
    title="Generate Token"
    :okText="steps[step].okText"
    @ok="steps[step].ok"
    @cancel="() => $emit('cancel')"
  >
    <a-row style="margin-bottom: 20px">
      <a-steps :current="step" size="small">
        <a-step title="Login">
          <a-icon slot="icon" type="user" />
        </a-step>
        <a-step title="Verification">
          <a-icon slot="icon" :type="step == 1 ? 'loading' : 'small-dash'" />
        </a-step>
        <a-step title="Save Token">
          <a-icon slot="icon" type="copy" />
        </a-step>
      </a-steps>
    </a-row>
    <template v-if="step == 0">
      <a-row>
        <a-alert v-if="error" :message="error" type="error" />
      </a-row>
      <a-form-model :model="credentials">
        <a-form-model-item>
          <a-input v-model="credentials.username" placeholder="Username">
            <a-icon
              slot="prefix"
              type="user"
              style="color: rgba(0, 0, 0, 0.25)"
            />
          </a-input>
        </a-form-model-item>
        <a-form-model-item>
          <a-input
            v-model="credentials.password"
            type="password"
            placeholder="Password"
          >
            <a-icon
              slot="prefix"
              type="lock"
              style="color: rgba(0, 0, 0, 0.25)"
            />
          </a-input>
        </a-form-model-item>
      </a-form-model>
    </template>
    <template v-else-if="step == 1">
      <a-row type="flex" justify="center">
        <a-col><h3>Requesting token...</h3></a-col>
      </a-row>
    </template>
    <template v-else-if="step == 2">
      <a-row>
        <h3>Copy Token</h3>
      </a-row>
      <a-row>
        <a-input :value="token" disabled>
          <a-icon slot="addonAfter" type="copy" @click="copyToken" />
        </a-input>
      </a-row>
    </template>
  </a-modal>
</template>

<script>
import { mapGetters } from "vuex";
import Clipboard from "@/mixins/clipboard";

export default {
  mixins: [Clipboard],
  props: {
    active: {
      required: true,
    },
  },
  computed: {
    ...mapGetters({ user: "loggedInUser" }),
  },
  data() {
    return {
      error: false,
      step: 0,
      token: false,
      credentials: {
        username: "",
        password: "",
      },
      steps: [
        {
          okText: "Submit",
          ok: this.authorize,
        },
        {
          okText: "...",
          ok: () => {},
        },
        {
          okText: "Done",
          ok: () => this.$emit("cancel"),
        },
      ],
    };
  },
  mounted() {
    this.credentials.username = this.user.username;
    this.$on("cancel", () => {
      this.token = "";
      this.error = false;
      this.step = 0;
      this.credentials = {
        username: "",
        password: "",
      };
    });
  },
  methods: {
    copyToken() {
      this.copyTextToClipboard(this.token);
      this.$message.success("Token copied to clipboard");
    },
    authorize() {
      this.step++;
      this.$axios({
        method: "post",
        url: "/api/account/token",
        data: this.credentials,
      })
        .then((res) => {
          this.token = res.data.token;
          this.step++;
        })
        .catch((e) => {
          this.error = e.response.data.message;
          this.step--;
        });
    },
  },
};
</script>