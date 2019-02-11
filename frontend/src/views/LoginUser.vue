<template>
  <v-container>
    <v-layout
      row
      wrap
      justify-center
      align-center
    >
      <v-card
        min-width="50%"
      >
        <v-card-title>
          <h1 class="mb-3">Login</h1>
        </v-card-title>
        <v-card-text>
          <v-text-field
            label="Username"
            v-model="userName"
          ></v-text-field>
          <v-text-field
            label="Password"
            v-model="password"
            :append-icon="show1 ? 'visibility_off' : 'visibility'"
            :type="show ? 'text' : 'password'"
            @click:append="show = !show"
          ></v-text-field>
          <p
            v-if="notLoggedIn"
            style="color: red"
          >
            Login required
          </p>
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
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      show: false,
      userName: "joe",
      password: "test123",
      notLoggedIn: false,
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
            localStorage.loginError = "";
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
      this.notLoggedIn = true;
    }
  }
};
</script>

<style scoped>
</style>
