<template>
  <v-container>
    <h1 class="mb-3">Device Management</h1>
    <v-card>
      <v-layout
        row wrap
      >
        <v-flex>
          <v-card>
            <v-card-title primary-title>
              <h2>Device hierarchy</h2>
            </v-card-title>
            <v-treeview
              v-model="tree"
              :items="items"
              activatable
              active-class="grey lighten-4 indigo--text"
              selected-color="indigo"
              open-on-click
              selectable
              expand-icon="mdi-chevron-down"
              on-icon="mdi-bookmark"
              off-icon="mdi-bookmark-outline"
              indeterminate-icon="mdi-bookmark-minus"
            >
          </v-treeview>
          </v-card>
        </v-flex>
        <v-flex>
          <v-card>
            <v-card-text>
              {{ tree[0] }}
            </v-card-text>
            <v-card-actions>
              <v-btn
                @click="newLevel"
              >
                Include new level
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-flex>
        <v-flex>
          <v-card>
            <v-card-text>
              {{ items }}
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
      counter: 0,
      items: [],
      tree: [],
      data: {
        "objects": [{
            "uid": "0x1119d",
            "name": "Johannes' Home",
            "objects": [{
                "uid": "0x1119e",
                "name": "First Floor",
                "objects": [{
                  "uid": "0x1119f",
                  "name": "Living Room",
                  "devices": [{
                    "uid": "0x111a0",
                    "name": "PC"
                  }]
                }]
              },
              {
                "uid": "0x111a3",
                "name": "Second Floor"
              }
            ],
            "devices": [{
              "uid": "0x111a4",
              "name": "le lamp"
            }]
          },
          {
            "uid": "0x111a5",
            "name": "Enclosing Room",
            "devices": [{
              "uid": "0x111a6",
              "name": "Enclosing-room-device"
            }]
          }
        ],
        "devices": [{
          "uid": "0x111a2",
          "name": "some device"
        }]
      }
    }
  },
  computed: {
    // objectTree() {
    //   return JSON.parse(this.realTree)
    // }
  },
  methods: {
    newLevel() {
      console.log("include new level")
    },
    addChildDevice(input, id) {
      for (let element of input) {
        if (element.id===id) {
          element.children.push({id: "testId", name: "testDevice"});
          return;
        }
        else if (element.children) {
          this.addChildDevice(element.children, id);
        }
      }
    },
    transformObject(input) {
      let res = {};

      res.id = input.uid;
      res.name = input.name;
      res.children = [];
      if (input.devices) {
        for (let device of input.devices) {
          res.children.push(this.transformObject(device));
        }
      }
      if (input.objects) {
        for (let object of input.objects) {
          res.children.push(this.transformObject(object));
        }
      }
      return res;
    },
   transform(input) {
      let res = [];

      for (let value of input.objects) {
        let el = this.transformObject(value);
        res.push(el)
      }
      for (let value of input.devices) {
        let el = this.transformObject(value);
        res.push(el)
      }
      return res;
    }
  },
  mounted() {
    this.items = this.transform(this.data);
    console.log(this.addChildDevice(this.items, "0x111a4"));
  }
}
</script>

<style>
</style>
