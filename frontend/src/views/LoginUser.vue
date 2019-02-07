<template>
  <v-container>
    <h1 class="mb-3">{{ headline }}</h1>
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
      errorMessage: "Login error",
      headline: "Please log in"
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
            localStorage.loginError = false;
            this.$router.push("/devices");
          }
        })
        .catch(e => {
          console.log(e);
          this.loginError = true;
        });
    }
  },
  created() {
    if (localStorage.loginError) {
      this.headline = "You must first log in to use the application";
    }
  }
};
</script>

<style scoped>
</style>
