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
            <v-layout column wrap>
              <v-flex
                mb-3
              >
                <div class="caption grey--text mb-1">
                  Active
                </div>
                <div>
                  {{ device.activated }}
                </div>
              </v-flex>
              <v-flex
                mb-3
              >
                <div class="caption grey--text mb-1">
                  Id
                </div>
                <div>
                  {{ device.id }}
                </div>
              </v-flex>
              <v-flex
                mb-3
              >
                <div class="caption grey--text mb-1">
                  Name
                </div>
                <div>
                  {{ device.name }}
                </div>
              </v-flex>
              <v-flex
                mb-3
              >
                <div class="caption grey--text mb-1">
                  Description
                </div>
                <div>
                  {{ device.description }}
                </div>
              </v-flex>
              <v-flex
                mb-3
              >
                <div class="caption grey--text mb-1">
                  Tags
                </div>
                <v-layout row wrap>
                  <div
                    v-for="tag in device.tags"
                    :key="tag"
                  >
                    <v-chip>{{ tag }}</v-chip>
                  </div>
                </v-layout>
              </v-flex>
            </v-layout>
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
      device: this.$store.getters.getDevice(this.$route.params.id),
      headers: ["Active", "Id", "Name", "Location", "Tags"]
    };
  }
};
</script>

<style lang="css">
</style>
