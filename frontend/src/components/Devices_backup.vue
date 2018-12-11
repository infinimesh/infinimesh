<template>
  <v-container>
    <v-layout>
      </v-toolbar>
        <v-card-text>
          <v-card-title>
            Devices
              <v-spacer></v-spacer>
              <v-text-field
                :append-icon="search?'close':'search'"
                :append-icon-cb="() => (search = '')"
                @keydown.native.escape="search=''"
                label="Search"
                single-line
                hide-details
                v-model="search"
              ></v-text-field>
          </v-card-title>
          <div>
            <v-card>

            </v-card>
            <v-data-table
              v-bind:headers="headers"
              v-bind:items="items"
              v-bind:search="search"
              v-bind:pagination.sync="pagination"
              hide-actions
              class="elevation-1"
            >
              <template slot="headerCell" scope="props">
                <v-tooltip bottom>
                  <span slot="activator">
                    {{ props.header.text }}
                  </span>
                  <span>
                    {{ props.header.text }}
                  </span>
                </v-tooltip>
              </template>
              <template slot="items" scope="props">
                <td>{{ props.item.status }}</td>
                <td  class="text-xs-right">{{ props.item.name }}</td>
                <td  class="text-xs-right">{{ props.item.location }}</td>
                <td  class="text-xs-right">{{ props.item.tags }}</td>
              </template>
            </v-data-table>
          </div>
          <div class="text-xs-center pt-2">
            <v-pagination v-model="pagination.page" :length="pages"></v-pagination>
          </div>
        </v-card-text>
        <div style="flex: 1 1 auto;"></div>
      </v-card>
    </v-dialog>
  </v-card>
</v-layout>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      dialog: false,
      fullScreen: true,
      search: "",
      pagination: {},
      selected: [],
      headers: [
        {
          text: "Status",
          align: "left",
          value: "name"
        },
        { text: "Name", value: "name" },
        { text: "Location", value: "location" },
        { text: "Tags", value: "tags" }
      ],
      items: [
        {
          value: false,
          status: "active",
          name: "Test device 1",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          status: "active",
          name: "Test device 2",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          name: "Test device 3",
          status: "active",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          name: "Test device 4",
          status: "active",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          name: "Test device 5",
          status: "active",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          name: "Test device 6",
          status: "active",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          status: "active",
          name: "Test device 7",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          status: "active",
          name: "Test device 8",
          location: "Düsseldorf",
          tags: "test"
        },
        {
          value: false,
          status: "inactive",
          name: "Test device 9",
          location: "Essen",
          tags: "test"
        },
        {
          value: false,
          status: "inactive",
          name: "Test device 10",
          location: "Essen",
          tags: "test"
        }
      ]
    };
  },
  computed: {
    pages() {
      return this.pagination.rowsPerPage
        ? Math.ceil(this.items.length / this.pagination.rowsPerPage)
        : 0;
    }
  }
};
</script>

<style>
</style>
