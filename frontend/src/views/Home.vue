<template>
  <v-container>
    <v-layout row wrap>
      <h1 class="mb-3">Your devices</h1>
      <devices></devices>
    </v-layout>
  </v-container>
</template>

<script>
import Devices from "../components/Devices.vue";

  export default {
    components: {
      Devices
    }
  }
</script>
