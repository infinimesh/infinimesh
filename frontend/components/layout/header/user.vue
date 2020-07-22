<template>
  <a-dropdown>
    <a-row type="flex" justify="end">
      <a-col :span="8">
        <a-avatar>{{ user.name[0].toUpperCase() }}</a-avatar>
      </a-col>
      <a-col class="user" :xs="{ span: 0 }" :sm="{ span: 0 }" :md="{ span: 16 }">
        {{ user.name }}
        <a-icon type="down" />
      </a-col>
    </a-row>
    <a-menu slot="overlay">
      <a-menu-item>
        <a href="javascript:;" @click="$store.commit('logout')">Log Out</a>
      </a-menu-item>
    </a-menu>
  </a-dropdown>
</template>

<script>
import Vue from "vue";
export default Vue.component("user-actions", {
  computed: {
    user() {
      return this.$store.getters.loggedInUser;
    }
  }
});
</script>