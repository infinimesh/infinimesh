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
              :items="items"
              :value="testValue"
              activatable
              active-class="grey lighten-4 indigo--text"
              selected-color="indigo"
              open-on-click
              selectable
              expand-icon="mdi-chevron-down"
              on-icon="mdi-checkbox-marked"
              off-icon="mdi-checkbox-blank-outline"
              indeterminate-icon="mdi-checkbox-blank-outline"
            >
          </v-treeview>
          </v-card>
        </v-flex>
        <v-flex>
          <v-card>
            <v-card-text>
              {{ topNode }}
            </v-card-text>
            <v-card-actions>
              <v-text-field
                label="Device name"
                clearable
                v-model="device.name"
              ></v-text-field>
              <v-btn
                @click="addNewLevel()"
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
      testValue: ["0x111a0"],
      device: {
        name: "",
        id: "",
        children: []
      },
      items: [],
      deviceTree: [],
      data: {
        objects: [
          {
            uid: "0x1119d",
            name: "Johannes' Home",
            objects: [
              {
                uid: "0x1119e",
                name: "First Floor",
                objects: [
                  {
                    uid: "0x1119f",
                    name: "Living Room",
                    devices: [
                      {
                        uid: "0x111a0",
                        name: "PC"
                      }
                    ]
                  }
                ]
              },
              {
                uid: "0x111a3",
                name: "Second Floor"
              }
            ],
            devices: [
              {
                uid: "0x111a4",
                name: "le lamp"
              }
            ]
          },
          {
            uid: "0x111a5",
            name: "Enclosing Room",
            devices: [
              {
                uid: "0x111a6",
                name: "Enclosing-room-device"
              }
            ]
          }
        ],
        devices: [
          {
            uid: "0x111a2",
            name: "some device"
          }
        ]
      }
    };
  },
  computed: {
    topNode() {
      return this.deviceTree[0];
    }
    // objectTree() {
    //   return JSON.parse(this.realTree)
    // }
  },
  methods: {
    selectTopNode() {
      return this.deviceTree[0];
    },
    addNewLevel() {
      this.device.id = Math.random().toString();
      let newDevice = JSON.parse(JSON.stringify(this.device));
      this.addChildDevice(this.items, this.topNode, newDevice);
      this.device.name = "";
      this.device.id = "";
      this.device.children = [];
    },
    addChildDevice(input, id, device) {
      for (let element of input) {
        if (element.id === id) {
          let newArr = element.children;
          newArr.push(device);
          element.children = newArr;
          return console.log("return");
        } else {
          if (element.children) {
            this.addChildDevice(element.children, id, device);
          }
          // little bug here: the function never enters the else loop below
          else {
            console.log("not found");
            return;
          }
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
        res.push(el);
      }
      for (let value of input.devices) {
        let el = this.transformObject(value);
        res.push(el);
      }
      return res;
    }
  },
  mounted() {
    this.items = this.transform(this.data);
  }
};
</script>

<style>
</style>
