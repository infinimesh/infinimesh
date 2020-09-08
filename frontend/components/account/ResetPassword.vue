<template>
  <a-modal
    :visible="active"
    :title="`Reset ${account.name} password`"
    okText="Reset"
    @ok="handleSubmit"
    @cancel="() => $emit('cancel')"
  >
    <a-form-model :model="model" :rules="rules" ref="resetAccountPasswordForm">
      <a-form-model-item prop="password" label="Password">
        <a-input-password placeholder="Enter password" v-model="model.password" />
      </a-form-model-item>
      <a-form-model-item prop="confirm_password" label="Confirm">
        <a-input-password placeholder="Confirm password" v-model="model.confirm_password" />
      </a-form-model-item>
    </a-form-model>
  </a-modal>
</template>

<script>
export default {
  props: {
    active: {
      required: true,
    },
    account: {
      required: true,
    },
  },
  data() {
    return {
      model: {},
      rules: {
        password: [
          { required: true, message: "Fill in the new password, please" },
        ],
        confirm_password: [
          { required: true, message: "Please, confirm the password" },
          {
            validator: (rule, value, raise) => {
              if (value != this.model.password) raise("oh noo");
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
      this.model = {
        password: "",
        confirm_password: "",
      };
    },
    handleSubmit() {
      let form = this.$refs.resetAccountPasswordForm;
      let errors = [];
      form.validateField(Object.keys(this.model), (err) => {
        if (err) {
          errors.push(err);
        }
      });
      if (errors.length === 0) this.$emit("reset", this.model.password);
    },
  },
};
</script>