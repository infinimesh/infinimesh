<template>
  <a-row type="flex" justify="center" class="rootRow" align="middle">
    <a-col :xs="22" :md="16" :lg="12">
      <a-row type="flex" justify="center">
        <h1>infinimesh Login</h1>
      </a-row>
      <a-row type="flex" justify="center"
        >Welcome to infinimesh. Log in with your username and password.</a-row
      >

      <a-row style="margin-top: 1rem">
        <a-form :form="form" @submit="handleSubmit">
          <a-form-item>
            <a-input
              v-decorator="[
                'username',
                {
                  rules: [
                    { required: true, message: 'Please input your username!' }
                  ]
                }
              ]"
              placeholder="Username"
            >
              <a-icon
                slot="prefix"
                type="user"
                style="color: rgba(0,0,0,.25)"
              />
            </a-input>
          </a-form-item>
          <a-form-item>
            <a-input
              v-decorator="[
                'password',
                {
                  rules: [
                    { required: true, message: 'Please input your Password!' }
                  ]
                }
              ]"
              type="password"
              placeholder="Password"
            >
              <a-icon
                slot="prefix"
                type="lock"
                style="color: rgba(0,0,0,.25)"
              />
            </a-input>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit" style="width: 100%"
              >Login</a-button
            >
          </a-form-item>
        </a-form>
      </a-row>
    </a-col>
  </a-row>
</template>

<script>
export default {
  data() {
    return {
      form: this.$form.createForm(this, { name: "login" })
    };
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault();
      this.form.validateFields(async (err, values) => {
        if (!err) {
          try {
            let res = await this.$auth.loginWith("local", {
              data: values
            });
            this.$router.push("/dashboard");
          } catch (e) {
            this.$notification.error({
              placement: "bottomLeft",
              duration: 10,
              ...e.response.data
            });
          }
        }
      });
    }
  }
};
</script>

<style scoped>
.rootRow {
  height: 100%;
  min-height: 500px;
  width: 100%;
}
</style>
