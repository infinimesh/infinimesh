<template>
  <v-app>
    <v-container>
      <v-row justify="center">
        <v-col sm="6" md="8" lg="6" cols="24">
          <v-row justify="center">
            <h1>infinimesh</h1>
          </v-row>
          <v-row justify="center" style="color: rgba(0, 0, 0, 0.65);"
            >Welcome to infinimesh. Log in with your username and
            password.</v-row
          >

          <v-row style="margin-top: 1rem">
            <a-form :form="form" @submit="handleSubmit" style="width: 100%">
              <a-form-item>
                <a-input
                  v-decorator="[
                    'username',
                    {
                      rules: [
                        {
                          required: true,
                          message: 'Please input your username!'
                        }
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
                        {
                          required: true,
                          message: 'Please input your Password!'
                        }
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
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </v-app>
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
            this.$router.push("/dashboard/devices");
          } catch (e) {
            this.$notification.error({
              placement: "bottomRight",
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

<style lang="less" scoped>
h1 {
  color: @line-color;
  font-family: Exo;
}
</style>
