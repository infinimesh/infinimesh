<template>
  <v-card
    flat
    width="50%"
  >
    <v-form
      v-model="form"
    >
    <v-text-field
      label="Username"
      v-model="user.name"
      :rules="rules.nameRules"
      clearable
    ></v-text-field>
    <v-layout
      row
      wrap
    >
      <v-text-field
        style="width: 40%"
        label="Password"
        v-model="user.passwordOne"
        :rules="rules.pwdOneRules"
        :append-icon="show ? 'visibility_off' : 'visibility'"
        :type="show ? 'text' : 'password'"
        @click:append="show = !show"
      ></v-text-field>
      <v-text-field
        style="width: 40%"
        label="Confirm password"
        class="ml-4"
        v-model="user.passwordTwo"
        :rules="rules.pwdTwoRules"
        :append-icon="show ? 'visibility_off' : 'visibility'"
        :type="show ? 'text' : 'password'"
        @click:append="show = !show"
      ></v-text-field>
    </v-layout>
  </v-form>
  </v-card>
</template>

<script>
export default {
  data() {
    return {
      show: false,
      form: false,
      user: {
        name: "",
        passwordOne: "",
        passwordTwo: ""
      },
      rules: {
        nameRules: [
          v => !!v || "Name is required",
          v =>
            !v.match(/(?=\W)(?=[^-_])/g) ||
            "Only alphanumerical characters and - _ allowed",
          v =>
            !this.$store.getters.getDevice(v) || "This device Id already exists"
        ],
        pwdOneRules: [
          v => !!v || "Password is required",
          v =>
            v.match(
              /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{10,})/g
            ) ||
            "Min length 10 char. Must have: 1 uppercase and 1 lowercase letter, 1 numeric char, 1 special char."
        ],
        pwdTwoRules: [
          v => (!!v && v) === this.user.passwordOne || "Passwords do not match"
        ]
      }
    };
  },
  watch: {
    form() {
      this.$emit("newUser", this.user);
    }
  }
};
</script>

<style scoped>
</style>
