<template>
  <div>
    <v-card-title>
      <h1>Device Registry</h1>
    </v-card-title>
    <v-card-text>
      <v-card flat>
        <v-card-title>
          <v-layout row wrap>
            <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details></v-text-field>
            <v-spacer></v-spacer>
            <v-spacer></v-spacer>
            <v-btn color="primary lighten-1" round :to="{ name: 'Register Device' }">
              <v-icon>add</v-icon>
            </v-btn>
          </v-layout>
        </v-card-title>
        <v-card-text>
          <v-layout row wrap>
            <v-flex id="scrollableCard">
              <v-data-table :headers="headers" :items="devices" :search="search" item-key="name" hide-actions>
                <template slot="items" slot-scope="props">
                                            <tr>
                                              <td class="text-xs-left" style="width: 10px">
                                                <v-icon v-if="props.item.enabled" color="green">check_circle</v-icon>
                                                <v-icon v-else color="grey">block</v-icon>
                                              </td>
                                              <td
                                                class="text-xs-left"
                                                style="cursor: pointer"
                                                @click="navigateTo(props.item.id)"
                                              >{{ props.item.id }}</td>
                                              <td v-if="props.item.tags" class="text-xs-left">
                                                <v-chip v-for="tag in props.item.tags" :key="tag">{{ tag }}</v-chip>
                                              </td>
                                              <td v-else class="text-xs-left">
                                              </td>
                                              <td class="text-xs-center">
                                                <v-menu offset-y>
                                                  <v-btn slot="activator" color="primary" flat>
                                                    <v-icon>more_vert</v-icon>
                                                  </v-btn>
                                                  <v-list>
                                                    <v-list-tile
                                                      v-for="option in options"
                                                      :key="option"
                                                      :to="{
                                                                      name: option,
                                                                      params: { id: props.item.id }
                                                                    }"
                                                    >
                                                      <v-list-tile-title>{{ option }}</v-list-tile-title>
                                                    </v-list-tile>
                                                  </v-list>
                                                </v-menu>
                                              </td> 
                                            </tr>
</template>
              </v-data-table>
            </v-flex>
          </v-layout>
        </v-card-text>
      </v-card>
    </v-card-text>
  </div>
</template>


<script>
import { mapGetters } from "vuex";

import Vue from "vue";

export default {
  data() {
    return {
      crumb: [
        {
          text: "Dashboard",
          disabled: false,
          href: "breadcrumbs_dashboard"
        },
        {
          text: "Link 1",
          disabled: false,
          href: "breadcrumbs_link_1"
        },
        {
          text: "Link 2",
          disabled: true,
          href: "breadcrumbs_link_2"
        }
      ],
      desserts: [
        {
          name: "Frozen Yogurt",
          calories: 159,
          fat: 6.0,
          carbs: 24,
          protein: 4.0,
          iron: "1%"
        }
      ],
      devicez: [
        {
          name: "x"
        }
      ],
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
      options: ["Unregister Device"]
    };
  },

  asyncComputed: {
    devices: {
      default() {
        return [];
      },
      get() {
        return new Promise((resolve, reject) => {
          return Vue.http
            .get(`namespaces/${this.namespace}/devices`)
            .then(res => {
              resolve(res.body.devices);
            })
            .catch(err => {
              reject(err);
            });

          // resolve([{
          // "name": "c"
          // }]);
        });
      }
      //         // if (!this.namespace) {
      //         //   resolve([])
      //         // }
      //         //   return Vue.http
      //         //   .get(`namespaces/joe/devices`)
      //         //     .then(res => res.json())
      //         //     .then(res => {
      //         //       // commit("storeNamespaces", res.namespaces);
      //         //       // commit("setNamespace", res.namespaces[1].name);
      //         //       resolve(res.devices);
      //         //     })
      //         //     .catch(err => {
      //         //       reject(err);
      //         //     });
      //         // });
      // Vue.http
      //           .get(`namespaces/${this.namespace}/devices`)
      //           .then(response => {
      //             response.body.devices;
      //             console.log("got it", response.body.devices)
      //             // this.bla = response.body.devices;
      //           })
      //           .catch(error => {
      //             console.log(error);
      //           });
      //         }

      // }
    }
  },
  computed: {
    pages() {
      return this.pagination.rowsPerPage
        ? Math.ceil(this.items.length / this.pagination.rowsPerPage)
        : 0;
    },
    // devices: {
    //   get: function() {

    // // let abc = await this.fetch()
    //     return [{"name": "d"}];
    //   }
    // },
    // return [];
    // console.log("xxx");
    // let abc = await this.fetch()
    // return function() {
    //   console.log("got" + abc);
    //   return abc;
    // }
    // },
    // set: function(chipsValues) {
    //   Vue.http
    //     .get(`namespaces/${this.namespace}/devices`)
    //     .then(response => {
    //       console.log(response.body.devices)
    //       return response.body.devices;
    //     })
    //     .catch(error => {
    //       console.log(error);
    //     }).get();

    //   console.log("set");
    //   // this.obj.values = chipsValues.map(val => {
    //   //   return {
    //   //     'name': val
    //   //   }
    //   // })
    //   if (!this.namespace) {
    //     return;
    //   }
    //   this.$store.dispatch("fetchDevices", this.namespace);

    // }
    // console.log("call");
    // // return this.$store.getters.getAllDevices;
    // if (!this.namespace) {
    //   return;
    // }
    // this.$store.dispatch("fetchDevices", this.namespace);
    // return this.$store.getters.getAllDevices;
    // // Vue.http
    //   .get(`namespaces/${this.namespace}/devices`)
    //   .then(response => {
    //     console.log(response.body.devices)
    //     this.abc = response.body.devices;
    //   })
    //   .catch(error => {
    //     console.log(error);
    //   });
    //   return [{
    //     id: "abc"
    //   }];
    // },
    ...mapGetters({
      namespace: "getNamespace"
    })
  },
  methods: {
    async fetch() {
      return new Promise((resolve, reject) => {
        return Vue.http
          .get(`namespaces/${this.namespace}/devices`)
          .then(response => {
            console.log(response.body.devices);
            return response.body.devices;
          })
          .catch(error => {
            console.log(error);
          })
          .get();
      });
    },
    navigateTo(id) {
      this.$router.push("devices/show/" + id);
    }
  },
  created() {
    console.log("namespace called in route", this.namespace);
  }
};
</script>

<style lang="css" scoped>
  #scrollableCard {
    max-height: 500px;
    overflow-y: auto;
  }
</style>
