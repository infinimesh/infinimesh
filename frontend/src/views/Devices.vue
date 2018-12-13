<template>
  <v-container xs9>
    <h1 class="mb-3">Your devices</h1>
    <v-card>
     <v-card-title>
       <v-text-field
         v-model="search"
         append-icon="search"
         label="Search"
         single-line
         hide-details
       ></v-text-field>
       <v-spacer></v-spacer>
       <v-spacer></v-spacer>
     </v-card-title>
     <v-data-table
       :headers="headers"
       :items="devices"
       :search="search"
       item-key="name"
     >
       <template slot="items" slot-scope="props">
         <tr>
           <td
           v-for="(attribute, i) in props.item"
           class="text-xs-left"
           style="cursor: pointer"
           :key="i"
            @click="navigateTo(props.item.deviceId)"
           >
           {{ attribute }}
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
              <v-icon>more_vert</v-icon>
             </v-btn>
             <v-list>
               <v-list-tile
                 v-for="(option, index) in options"
                 :key="index"
                 :to="{name: option, params: { id: props.item.deviceId }}"
               >
                 <v-list-tile-title>{{ option }}</v-list-tile-title>
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
     <v-card-actions>
       <v-btn
        color="primary lighten-1"
        bottom
        left
        round
        :to="{ name: 'Register device' }"
        >
        <v-icon>add</v-icon>
      </v-btn>
     </v-card-actions>
   </v-card>
 </v-container>
</template>

<script>
export default {
  data() {
    return {
      selected: true,
      search: "",
      headers: [
        {
          text: "Status",
          align: "left",
          value: "status"
        },
        {
          text: "Id",
          align: "left",
          value: "deviceId"
        },
        {
          text: "Name",
          align: "left",
          value: "name"
        },
        {
          text: "Location",
          align: "left",
          value: "location"
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
      options: ["Delete device"]
      // devices: [
      //   {
      //     status: "active",
      //     id: 1,
      //     name: "Device 1",
      //     location: "Düsseldorf",
      //     tags: "test"
      //   },
      //   {
      //     status: "active",
      //     id: 2,
      //     name: "Device 2",
      //     location: "Essen",
      //     tags: "test"
      //   },
      //   {
      //     status: "active",
      //     id: 3,
      //     name: "Device 3",
      //     location: "Berlin",
      //     tags: "test"
      //   },
      //   {
      //     status: "active",
      //     id: 4,
      //     name: "Device 4",
      //     location: "Düsseldorf",
      //     tags: "prod"
      //   },
      //   {
      //     status: "active",
      //     id: 5,
      //     name: "Device 5",
      //     location: "Düsseldorf",
      //     tags: "test"
      //   },
      //   {
      //     status: "inactive",
      //     id: 6,
      //     name: "Device 6",
      //     location: "Düsseldorf",
      //     tags: "test"
      //   },
      //   {
      //     status: "active",
      //     id: 7,
      //     name: "Device 7",
      //     location: "Essen",
      //     tags: "test"
      //   }
      // ]
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
      this.$router.push("devices/" + id);
    }
  }
};
</script>

<style>
</style>
