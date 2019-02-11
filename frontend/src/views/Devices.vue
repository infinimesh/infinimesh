<template>
  <div>
    <v-card-title>
      <h1>
        Device registry
      </h1>
    </v-card-title>
    <v-card-text>
      <v-card
        flat
      >
        <v-card-title>
          <v-layout row wrap>
            <v-text-field
               v-model="search"
               append-icon="search"
               label="Search"
               single-line
               hide-details
            >
            </v-text-field>
            <v-spacer></v-spacer>
            <v-spacer></v-spacer>
            <v-btn
              color="primary lighten-1"
              round
              :to="{ name: 'Register device' }"
            >
              <v-icon>
                add
              </v-icon>
            </v-btn>
          </v-layout>
        </v-card-title>
        <v-card-text>
          <v-layout row wrap>
            <v-flex
              id="scrollableCard"
            >
              <v-data-table
               :headers="headers"
               :items="devices"
               :search="search"
               item-key="name"
               hide-actions
              >
                <template slot="items" slot-scope="props">
                  <tr>
                    <td
                      class="text-xs-left"
                      style="width: 10px"
                    >
                    <v-icon
                      v-if="props.item.enabled"
                      color="green"
                    >
                      check_circle
                   </v-icon>
                   <v-icon
                     v-else
                     color="grey"
                   >
                    block
                  </v-icon>
                    </td>
                    <td
                      class="text-xs-left"
                      style="cursor: pointer"
                      @click="navigateTo(props.item.id)"
                    >
                      {{ props.item.id }}
                    </td>
                    <td
                      class="text-xs-left"
                    >
                      <v-chip
                        v-for="tag in props.item.tags"
                        :key="tag"
                      >
                        {{ tag }}
                      </v-chip>
                    </td>
                    <td
                      class="text-xs-center"
                    >
                    <v-menu offset-y>
                      <v-btn
                        slot="activator"
                        color="primary"
                        flat
                      >
                        <v-icon>
                          more_vert
                        </v-icon>
                      </v-btn>
                      <v-list>
                        <v-list-tile
                          v-for="option in options"
                          :key="option"
                          :to="{name: option, params: { id: props.item.id }}"
                        >
                          <v-list-tile-title>
                            {{ option }}
                          </v-list-tile-title>
                        </v-list-tile>
                      </v-list>
                    </v-menu>
                  </td>
                </tr>
              </template>
              <v-alert slot="no-results" :value="true" color="error" icon="warning">
                Your search for "{{ search }}" found no results.
              </v-alert>
            </v-data-table>
            </v-flex>
          </v-layout>
        </v-card-text>
        </v-card>
    </v-card-text>
    <v-divider />
  </div>
</template>

<script>
export default {
  data() {
    return {
      selected: true,
      search: "",
      headers: [
        {
          text: "Active",
          align: "left",
          value: "enabled"
        },
        {
          text: "Id",
          align: "left",
          value: "id"
        },
        {
          text: "Tags",
          align: "left",
          value: "tags"
        },
        {
          text: "Further actions",
          align: "center"
        }
      ],
      options: ["Unregister device"]
    };
  },
  computed: {
    pages() {
      return this.pagination.rowsPerPage
        ? Math.ceil(this.items.length / this.pagination.rowsPerPage)
        : 0;
    },
    devices() {
      return this.$store.getters.getAllDevices;
    }
  },
  methods: {
    navigateTo(id) {
      this.$router.push("devices/show/" + id);
    }
  },
  created() {
    this.$store.dispatch("fetchDevices");
  }
};
</script>

<style lang="css" scoped>
  #scrollableCard {
  max-height: 500px;
  overflow-y: auto;
  }
</style>
