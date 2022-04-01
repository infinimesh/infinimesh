<template>
  <router-view />
</template>

<script>
export default {
  computed: {
    isLoggedIn() {
      return false
    }
  },
  created() {
    console.log(this)

  }
}
</script>