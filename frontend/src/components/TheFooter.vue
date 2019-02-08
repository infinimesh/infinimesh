<template>
  <div
    class="mt-5"
  >
    <v-footer
      height="auto"
      color="primary lighten-1"
      absolute
    >
      <v-layout
        justify-center
        row
        wrap
      >
        <v-btn
          v-for="link in links"
          :key="link"
          color="white"
          flat
          round
        >
          {{ link }}
        </v-btn>
      </v-layout>
    </v-footer>
</div>
</template>

<script>
export default {
  data() {
    return {
      links: ["HOME", "ABOUT US", "CONTACT US"]
    };
  }
};
</script>

<style>
</style>
