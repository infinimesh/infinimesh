<template>
  <v-container>
    <h1 class="mb-3">Device Management</h1>
    <v-card>
      <v-layout
        row wrap
      >
        <v-flex>
          <v-card
            flat
            class="ma-2"
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
              <template
                v-if="active"
                slot="append"
                slot-scope="{ item, active }"
              >
                <v-icon
                  :color="active ? 'primary' : ''"
                  @click.stop="showNodePanel=true"
                >
                  add
                </v-icon>
                <v-icon
                  :color="active ? 'primary' : ''"
                  @click.stop="deleteBranch"
                >
                  delete
                </v-icon>
              </template>
              <template
                slot="prepend"
                slot-scope="{ item }"
                v-if = "item.type === 'device'"
              >
                <v-icon>
                  smartphone
                </v-icon>
              </template>
            </v-treeview>
          </v-card>
        </v-flex>
        <v-spacer
          v-if="!showNodePanel"
        ></v-spacer>
        <v-divider
          v-if="showNodePanel"
          vertical
        ></v-divider>
        <v-flex
        >
          <v-card
            class="ma-2"
            flat
          >
            <div
              v-if="showNodePanel"
            >
              <v-layout
                align-end
                justify-end
              >
                <v-icon
                  style="cursor: pointer"
                  @click="showNodePanel = false"
                  class="ma-3"
                >
                  close
                </v-icon>
              </v-layout>
              <v-card-text>
                <v-text-field
                  label="Name of new node"
                  clearable
                  v-model="node.name"
                ></v-text-field>
                <v-text-field
                  label="Type of new node"
                  clearable
                  v-model="node.type"
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
                <v-alert
                  :value="alert.value"
                  type="warning"
                >
                  {{ alert.message }}
                </v-alert>
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
                  @click="revert"
                >
                  Revert
                </v-btn>
              </v-card-actions>
            </div>
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
        type: "",
        children: []
      },
      items: [],
      nodeHistory: [],
      nodeAdderFunction: "",
      radioLabels: ["Add child", "Add sibling", "Attach to new parent"],
      showNodePanel: false,
      alert: {
        value: false,
        message: ""
      }
    };
  },
  methods: {
    addNewNode() {
      this.alert.value = false;
      if (this.checkIfName()) {
        setTimeout(() => (this.alert.value = false), 2000);
        return;
      }
      this.extendNodeHistory();
      let newNode = this.setNode();
      switch (this.nodeAdderFunction) {
        case "Add child":
          this.addChildNode(this.items, this.active[0], newNode);
          break;
        case "Add sibling":
          this.addSiblingNode(this.items, this.active[0], newNode);
          break;
        case "Attach to new parent":
          this.attachToNewParentNode(this.items, this.active[0], newNode);
          break;
      }
      this.clearNode();
    },
    deleteBranch() {
      console.log("Delete Branch");
    },
    checkIfName() {
      if (!this.node.name) {
        this.alert.message = "Node must have a name";
        this.alert.value = true;
        return true;
      }
    },
    extendNodeHistory() {
      if (this.nodeHistory.length <= 5) {
        this.nodeHistory.push(JSON.parse(JSON.stringify(this.items)));
      }
    },
    setNode() {
      this.node.id = Math.random().toString();
      let newNode = JSON.parse(JSON.stringify(this.node));
      return newNode;
    },
    clearNode() {
      this.node.name = "";
      this.node.id = "";
      this.node.children = [];
    },
    addChildNode(input, id, node) {
      for (let element of input) {
        if (element.id === id) {
          let newArr = element.children;
          newArr.push(node);
          element.children = newArr;
          return node.id;
        } else if (element.children) {
          this.addChildNode(element.children, id, node);
        }
      }
    },
    addSiblingNode(input, id, node) {
      for (let element of input) {
        if (element.id === id) {
          input.splice(input.indexOf(element) + 1, 0, node);
          return node.id;
        } else if (element.children) {
          this.addSiblingNode(element.children, id, node);
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
          return node.id;
        } else if (element.children) {
          this.attachToNewParentNode(element.children, id, node);
        }
      }
    },
    transformObject(input) {
      let res = {};

      res.id = input.uid;
      res.name = input.name;
      res.type = input.type;
      res.children = [];
      if (input.devices) {
        for (let device of input.devices) {
          device.type = "device";
          res.children.push(this.transformObject(device));
        }
      }
      if (input.objects) {
        for (let object of input.objects) {
          object.type = "node";
          res.children.push(this.transformObject(object));
        }
      }
      return res;
    },
    transform(input) {
      let res = [];

      for (let value of input.objects) {
        value.type = "node";
        let el = this.transformObject(value);
        el.type = "node";
        res.push(el);
      }
      for (let value of input.devices) {
        value.type = "device";
        let el = this.transformObject(value);
        el.type = "device";
        res.push(el);
      }
      return res;
    },
    revert() {
      if (this.nodeHistory.length) {
        this.items = this.nodeHistory.pop();
      } else {
        this.alert.message = "Further reverts not possible";
        this.alert.value = true;
        setTimeout(() => (this.alert.value = false), 2000);
      }
    }
  },
  created() {
    this.$http.get("objects").then(response => {
      this.items = this.transform(response.body);
    });
  }
};
</script>

<style>
</style>
