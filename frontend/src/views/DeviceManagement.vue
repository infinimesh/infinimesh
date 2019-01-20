<template>
  <v-container>
    <h1 class="mb-3">Device Management</h1>
    <v-card>
      <v-layout
        row wrap
      >
        <v-flex>
          <v-card
          max-width="400"
          flat
          >
            <v-card-title primary-title>
              <h2>Device hierarchy</h2>
            </v-card-title>
            <v-treeview
              :items="items"
              activatable
              :active.sync = "active"
              active-class="grey lighten-4 indigo--text"
              selected-color="indigo"
            >
            <v-icon
              v-if="active"
              slot="append"
              slot-scope="{ item, active }"
              :color="active ? 'primary' : ''"
              @click.stop="showNodePanel=true"
            >
              add
            </v-icon>
          </v-treeview>
          </v-card>
        </v-flex>
        <v-flex
          v-if="showNodePanel"
        >
          <v-card
            class="ma-2"
            flat
          >
            <v-card-text>
              <v-text-field
                label="Name of new node"
                clearable
                v-model="node.name"
              ></v-text-field>
              <v-radio-group
                v-model="nodeAdderFunction">
                <v-radio
                  v-for="(label, i) in radioLabels"
                  :key="i"
                  :label="label"
                  :value="label"
                ></v-radio>
              </v-radio-group>
            </v-card-text>
              <v-card-actions>
              <v-btn
                round
                @click="addNewNode()"
                class="mr-3"
              >
                Include new level
              </v-btn>
              <v-btn
                round
              >
                Revert
              </v-btn>
            </v-card-actions>
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
      active: [],
      node: {
        name: "",
        id: "",
        children: []
      },
      items: [],
      nodeAdderFunction: "",
      radioLabels: ["Add child", "Add sibling", "Attach to new parent"],
      showNodePanel: false
    };
  },
  computed: {},
  methods: {
    addNewNode() {
      this.node.id = Math.random().toString();
      let newDevice = JSON.parse(JSON.stringify(this.node));
      switch (this.nodeAdderFunction) {
        case "Add child":
          this.addChildNode(this.items, this.active[0], newDevice);
          break;
        case "Add sibling":
          this.addSiblingNode(this.items, this.active[0], newDevice);
          break;
        case "Attach to new parent":
          this.attachToNewParentNode(this.items, this.active[0], newDevice);
          break;
      }
      this.node.name = "";
      this.node.id = "";
      this.node.children = [];
    },
    addChildNode(input, id, node) {
      console.log("input", input, "id", id, "node", node);
      for (let element of input) {
        if (element.id === id) {
          let newArr = element.children;
          newArr.push(node);
          element.children = newArr;
        } else {
          if (element.children) {
            this.addChildNode(element.children, id, node);
          }
          // little bug here: the function never enters the else loop below
          else {
            console.log("not found");
            return;
          }
        }
      }
    },
    addSiblingNode(input, id, node) {
      for (let element of input) {
        if (element.id === id) {
          input.splice(input.indexOf(element) + 1, 0, node);
          return node.id;
        } else {
          if (element.children) {
            this.addSiblingNode(element.children, id, node);
          }
          // little bug here: the function never enters the else loop below
          else {
            console.log("not found");
            return;
          }
          console.log("returns");
        }
      }
    },
    attachToNewParentNode(input, id, node) {
      for (let element of input) {
        if (element.id === id) {
          node.children.push(element);
          let newNode = JSON.parse(JSON.stringify(node));
          this.addSiblingNode(this.items, id, newNode);
          input.splice(input.indexOf(element), 1);
          return;
        } else if (element.children) {
          this.attachToNewParentNode(element.children, id, node);
        } else {
          return "Error";
        }
      }
    },
    transformObject(input) {
      let res = {};

      res.id = input.uid;
      res.name = input.name;
      res.children = [];
      if (input.devices) {
        for (let node of input.devices) {
          res.children.push(this.transformObject(node));
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
    this.$http.get("objects").then(response => {
      this.items = this.transform(response.body);
    });
  }
};
</script>

<style>
</style>
