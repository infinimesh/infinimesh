<template>
  <a-row type="flex" justify="center" class="rootRow" align="middle">
    <a-col :span="12">
      <a-row type="flex" justify="center">
        <h1>infinimesh Login</h1>
      </a-row>
      <a-row
        type="flex"
        justify="center"
      >Welcome to infinimesh. Log in with your username and password.</a-row>

      <a-row>
        <a-form>
          <a-form-item
            label="Username"
            v-decorator="['username', {rules: [{required: true, message: 'Pleeeeease, input your username, OK?!'}]}]"
          >
            <a-input placeholder="Enter your username" />
          </a-form-item>
          <a-form-item
            label="Password"
            v-decorator="['password', {rules: [{required: true, message: 'Pleeeeease, input your password, OK?!'}]}]"
          >
            <a-input-password placeholder="Input password" />
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
