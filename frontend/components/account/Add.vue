<template>
  <a-drawer
    title="Create New Account"
    :visible="active"
    :width="drawerSize"
    :headerStyle="{ fontSize: '4rem' }"
    @close="$emit('cancel')"
  >
    <a-form-model
      :model="account"
      :rules="rules"
      :label-col="{ xs: 24, sm: 6, md: 6, lg: 6 }"
      :wrapper-col="{ xs: 24, sm: 16, md: 18, lg: { span: 14, offset: 1 } }"
      ref="accountAddForm"
    >
      <a-form-model-item prop="name" label="Username">
        <a-input v-model="account.name" type="email" />
      </a-form-model-item>
      <a-form-model-item prop="password" label="Password">
        <a-input-password
          placeholder="Enter password"
          v-model="account.password"
        />
      </a-form-model-item>
      <a-form-model-item prop="confirm_password" label="Confirm">
        <a-input-password
          placeholder="Confirm password"
          v-model="account.confirm_password"
        />
      </a-form-model-item>
    </a-form-model>
    <div
      :style="{
        position: 'absolute',
        right: 0,
        bottom: 0,
        width: '100%',
        borderTop: '1px solid #e9e9e9',
        padding: '10px 16px',
        textAlign: 'right',
        zIndex: 1,
      }"
      id="accountAddDrawerActionsRow"
    >
      <a-button
        :style="{ marginRight: '8px' }"
        @click="setDefault"
        class="ant-btn-danger"
        >Reset</a-button
      >
      <a-button
        :style="{ marginRight: '8px' }"
        @click="$emit('cancel')"
        type="primary"
        >Cancel</a-button
      >
      <a-button type="success" @click="handleSubmit">Submit</a-button>
    </div>
  </a-drawer>
</template>

<script>
import Vue from "vue";

import drawerSizeMixin from "@/mixins/drawer-size.js";

export default Vue.component("account-add", {
  mixins: [drawerSizeMixin],
  props: {
    active: {
      required: true,
    },
  },
  data() {
    return {
      account: {},
      rules: {
        name: [
          {
            required: true,
            message: "Please, fill in the user email",
          },
          {
            type: "email",
            message: "Should be a valid email address",
          },
        ],
        password: [
          { required: true, message: "Fill in the new password, please" },
        ],
        confirm_password: [
          { required: true, message: "Please, confirm the password" },
          {
            validator: (rule, value, raise) => {
              if (value != this.account.password) raise("oh noo");
            },
            message: "Passwords don't match",
          },
        ],
      },
    };
  },
  watch: {
    active: "setDefault",
  },
  mounted() {
    this.setDefault();
  },
  methods: {
    setDefault() {
      this.account = {
        name: "",
        password: "",
        confirm_password: "",
      };
    },
    handleSubmit() {
      let form = this.$refs["accountAddForm"];
      let errors = [];
      form.validateField(Object.keys(this.account), (err) => {
        if (err) {
          errors.push(err);
        }
      });
      if (errors.length === 0) {
        delete this.account.confirm_password;
        this.$emit("add", { account: this.account });
      }
    },
  },
});
</script>
