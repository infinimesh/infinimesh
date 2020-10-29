<template>
  <v-app>
    <v-container style="margin-top: 100px">
      <a-row type="flex" justify="center">
        <a-col sm="6" md="8" lg="6">
          <a-row type="flex" justify="center">
            <h1>infinimesh</h1>
          </a-row>
          <a-row
            type="flex"
            justify="center"
            style="color: var(--logo-color); opacity: 0.65; text-align: center"
            >Welcome to infinimesh. Log in with your username and
            password.</a-row
          >

          <a-row style="margin-top: 1rem">
            <a-form :form="form" @submit="handleSubmit" style="width: 100%">
              <a-form-item>
                <a-input
                  v-decorator="[
                    'username',
                    {
                      rules: [
                        {
                          required: true,
                          message: 'Please input your username!',
                        },
                      ],
                    },
                  ]"
                  placeholder="Username"
                >
                  <a-icon
                    slot="prefix"
                    type="user"
                    style="color: rgba(0, 0, 0, 0.25)"
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
                          message: 'Please input your Password!',
                        },
                      ],
                    },
                  ]"
                  type="password"
                  placeholder="Password"
                >
                  <a-icon
                    slot="prefix"
                    type="lock"
                    style="color: rgba(0, 0, 0, 0.25)"
                  />
                </a-input>
                <a-progress
                  v-if="login_progress"
                  :percent="100"
                  status="active"
                  :show-info="false"
                  class="login-progress"
                ></a-progress>
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
    </v-container>
    <v-footer app>
      <infinimesh-footer />
    </v-footer>
  </v-app>
</template>

<script>
import InfinimeshFooter from "@/components/generic/footer.vue";

export default {
  components: {
    InfinimeshFooter,
  },
  data() {
    return {
      login_progress: false,
      form: this.$form.createForm(this, { name: "login" }),
    };
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault();
      this.login_progress = true;
      this.form.validateFields(async (err, values) => {
        if (!err) {
          try {
            let res = await this.$auth.loginWith("local", {
              data: values,
            });
            this.$router.push("/dashboard/devices");
          } catch (e) {
            this.$notification.error({
              duration: 10,
              ...e.response.data,
            });
          }
        }
        this.login_progress = false;
      });
    },
  },
};
</script>

<style scoped>
.rootRow {
  height: 100%;
  min-height: 500px;
  width: 100%;
}
h1 {
  color: var(--logo-color);
  font-family: Exo;
}
.login-progress {
  height: 18px;
}
</style>
<style>
.login-progress > .ant-progress-success-bg,
.ant-progress-bg {
  background-color: var(--primary-color);
}
</style>