<template>
  <v-container>
    <h1 class="mb-3">Login</h1>
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
      password: "test123"
    };
  },
  methods: {
    login() {
      this.$http.post("token", {
        username: this.userName,
        password: this.password
      })
      .then((response) => {
        if (response.status === 200) {
          localStorage.token = response.body.token;
          console.log(localStorage);
        }
      })
      .catch((e) => console.log(e))
    }
  },
  mounted() {
    console.log(localStorage.token);
  }
};
</script>

<style scoped>
</style>
