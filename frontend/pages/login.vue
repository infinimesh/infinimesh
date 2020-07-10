<template>
  <a-row type="flex" justify="center" class="rootRow" align="middle">
    <a-col :span="12">
      <a-row>
        <h1>Login</h1>
      </a-row>
      <a-row>
        <a-form>
          <a-form-item></a-form-item>
        </a-form>
      </a-row>
    </a-col>
  </a-row>
</template>

<script lang="ts">
import Vue from "vue";

interface IUser {
  login: string;
  password: string;
}

export default Vue.extend({
  data() {
    const form = {};
    return {
      form
    };
  }
});
</script>

<style scoped>
.rootRow {
  height: 100%;
  min-height: 500px;
  width: 100%;
}
</style>
