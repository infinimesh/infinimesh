<template>
  <v-container>
    <h1 class="mb-3">Successfully logged out. See you soon!</h1>
    <span>We redirect you to the login page in 10 seconds ... </span>
    <v-card>
    </v-card>
  </v-container>
</template>

<script>
export default {
  created() {
    setTimeout(() => this.$router.push({ path: "/login" }), 10000);
  }
};
</script>
<style scoped>
</style>
