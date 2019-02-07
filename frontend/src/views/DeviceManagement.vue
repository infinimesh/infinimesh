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
                  @click.stop="activeComp = 'addNode'"
                >
                  add
                </v-icon>
                <v-icon
                  :color="active ? 'primary' : ''"
                  @click.stop="activeComp = 'deleteNode'"
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
          v-if="!activeComp"
        ></v-spacer>
        <v-divider
          v-if="activeComp"
          vertical
        ></v-divider>
        <v-flex
        >
          <v-card
            class="ma-2"
            flat
          >
            <div
              v-if="activeComp"
            >
              <v-layout
                align-end
                justify-end
              >
                <v-icon
                  style="cursor: pointer"
                  @click="activeComp = ''"
                  class="ma-3"
                >
                  close
                </v-icon>
              </v-layout>
              <div
                v-if="activeComp === 'addNode'"
              >
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
              <div
                v-if="activeComp === 'deleteNode'"
              >
                <v-card-title primary-title>
                  Are you sure you want to delete this node?
                </v-card-title>
                <v-alert
                  :value="alert.value"
                  type="warning"
                >
                  {{ alert.message }}
                </v-alert>
                <v-card-actions>
                  <v-btn
                    round
                    @click="deleteNode()"
                    class="mr-3"
                  >
                    Confirm
                  </v-btn>
                </v-card-actions>
              </div>
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
      activeComp: "",
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
        parent: this.active[0],
        name: this.setNode().name
      };
      this.$store.dispatch("addChildNode", payload);
      this.clearNode();
    },
    deleteNode() {
      this.$store.dispatch("deleteNode", this.active[0]);
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
    }
  },
  created() {
    this.$store.dispatch("fetchNodeTree").catch(e => console.log(e));
  }
};
</script>

<style>
</style>
