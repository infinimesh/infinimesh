<template>
  <v-container>
    <h1 class="mb-3">Please log in</h1>
    <v-card>
      <v-card-text>
        <v-text-field
          label="Username"
          v-model="userName"
        ></v-text-field>
        <v-text-field
          label="Password"
          v-model="password"
        ></v-text-field>
      </v-card-text>
      <v-alert
       v-model="loginError"
       type="error"
       icon="error"
      >
        {{ errorMessage }}
      </v-alert>
      <v-card-actions>
        <v-btn
        color="primary"
        @click="login"
        >
          Login
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      userName: "joe",
      password: "test123",
      loginError: false,
      errorMessage: "Login error"
    };
  },
  methods: {
    login() {
      this.$http
        .post("token", {
          username: this.userName,
          password: this.password
        })
        .then(response => {
          if (response.status === 200) {
            localStorage.token = response.body.token;
            this.$router.push("/devices");
          }
        })
        .catch(e => {
          console.log(e);
          this.loginError = true;
        });
    }
  }
};
</script>

<style scoped>
</style>
