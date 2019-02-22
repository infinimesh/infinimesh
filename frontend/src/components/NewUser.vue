<template>
  <v-card flat width="50%">
    <v-form v-model="form">
      <v-text-field
        label="Username"
        v-model="user.name"
        box
        clearable
      ></v-text-field>
      <v-layout row wrap>
        <v-text-field
          label="Password"
          v-model="user.passwordOne"
          box
          class="mr-4"
          :append-icon="show ? 'visibility_off' : 'visibility'"
          :type="show ? 'text' : 'password'"
          @click:append="show = !show"
        ></v-text-field>
        <v-text-field
          label="Confirm password"
          v-model="user.passwordTwo"
          box
          :rules="pwdRules"
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
  data: () => ({
    show: false,
    form: false,
    test: "abc",
    user: {
      name: "",
      passwordOne: "",
      passwordTwo: ""
    },
    pwdRules: [
      v => (!!v && v) === this.user.passwordOne || "Passwords do not match"
    ]
  })
};
</script>

<style scoped></style>
