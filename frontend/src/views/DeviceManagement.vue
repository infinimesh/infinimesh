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
              :items="nodeTree"
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
                  v-if="item.type === 'node'"
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
      nodeAdderFunction: "",
      radioLabels: ["Add child", "Add sibling", "Attach to new parent"],
      showNodePanel: false,
      alert: {
        value: false,
        message: ""
      }
    };
  },
  computed: {
    nodeTree() {
      return this.$store.getters.getNodeTree;
    }
  },
  methods: {
    addNewNode() {
      this.alert.value = false;
      if (this.checkIfName()) {
        setTimeout(() => (this.alert.value = false), 2000);
        return;
      }
      let payload = {
        id: this.active[0],
        node: this.setNode()
      };
      this.$store.dispatch("addChildNode", payload);
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
    setNode() {
      this.node.id = Math.random().toString();
      let newNode = JSON.parse(JSON.stringify(this.node));
      return newNode;
    },
    clearNode() {
      this.node.name = "";
      this.node.id = "";
      this.node.type = "";
      this.node.children = [];
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
    }
  },
  created() {
    this.$store.dispatch("fetchNodeTree");
  }
};
</script>

<style>
</style>
