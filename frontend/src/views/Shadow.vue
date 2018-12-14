<template>
  <v-container>
    <h1 class="mb-3">Device overview</h1>
    <v-card>
      <v-layout row wrap>
        <v-flex>
          <v-card flat>
            <v-card-title primary-title>
              <h3>Status</h3>
            </v-card-title>
            <v-card-text>
              Desired state:
            </v-card-text>
            <v-card-text>
              Actual state:
            </v-card-text>
          </v-card>
        </v-flex>
        <v-flex>
          <v-card flat>
            <v-card-title primary-title>
              <h3>Device information</h3>
            </v-card-title>
            <v-card-text>
              Device Id: {{ device.deviceId }}
            </v-card-text>
            <v-card-text>
              Device name: {{ device.name }}
            </v-card-text>
            <v-card-text>
              Device description: {{ device.description }}
            </v-card-text>
            <v-card-text>
              Device location: {{ device.location }}
            </v-card-text>
            <v-card-text>
              Device tags: {{ device.tags }}
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      device: this.$store.getters.getDevice(parseInt(this.$route.params.id))
      }
    }
  }
</script>

<style lang="css">
</style>
