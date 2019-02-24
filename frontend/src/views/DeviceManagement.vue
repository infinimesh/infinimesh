<template>
  <div>
    <v-card-title>
      <h1 class="mb-3">Display, manage, modify and organize your devices</h1>
    </v-card-title>
    <v-divider></v-divider>
    <v-layout row wrap>
      <v-flex>
        <v-card flat class="ma-2" color="grey lighten-4">
          <v-sheet class="pa-3 primary lighten-2">
            <v-text-field
              v-model="search"
              label="Search device registry"
              dark
              flat
              solo-inverted
              hide-details
              clearable
              clear-icon="mdi-close-circle-outline"
            ></v-text-field>
          </v-sheet>
          <v-card-text>
            <v-treeview
              :items="nodeTree"
              :search="search"
              activatable
              :active.sync="active"
              active-class="grey lighten-4 indigo--text"
              selected-color="indigo"
            >
              <template slot="label" slot-scope="{ item }">
                <drag-drop-slot
                  :class="['tree-item', (over && over.id === item.id ? over.mode : '')]"
                  style="cursor: pointer"
                  :key="item.id"
                  :item="item"
                  @drag="drag"
                  @enter="enter"
                  @leave="leave"
                  @hover="hover"
                  @drop="drop"
                >
                  <v-icon>
                    drag_indicator
                  </v-icon>
                  <span>
                    {{ item.name }}
                  </span>
                </drag-drop-slot>
              </template>
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
                v-if="item.type === 'device'"
              >
                <v-icon>
                  smartphone
                </v-icon>
              </template>
            </v-treeview>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-divider vertical></v-divider>
      <v-spacer v-if="!activeComp"></v-spacer>
      <v-flex>
        <v-card class="ma-2" flat>
          <div v-if="activeComp">
            <v-layout align-end justify-end>
              <v-icon
                style="cursor: pointer"
                @click="activeComp = ''"
                class="ma-3"
              >
                close
              </v-icon>
            </v-layout>
            <div v-if="activeComp === 'addNode'">
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
                <v-alert :value="alert.value" type="warning">
                  {{ alert.message }}
                </v-alert>
              </v-card-text>
              <v-card-actions>
                <v-btn round @click="addNewNode()" class="mr-3">
                  Include new level
                </v-btn>
              </v-card-actions>
            </div>
            <div v-if="activeComp === 'deleteNode'">
              <v-card-title primary-title>
                Are you sure you want to delete this node?
              </v-card-title>
              <v-alert :value="alert.value" type="warning">
                {{ alert.message }}
              </v-alert>
              <v-card-actions>
                <v-btn round @click="deleteNode()" class="mr-3">
                  Confirm
                </v-btn>
              </v-card-actions>
            </div>
          </div>
        </v-card>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { DragDropContext } from "vue-react-dnd";
import HTML5Backend from "react-dnd-html5-backend";
import DragDropSlot from "../components/DragDropSlot.vue";

export default {
  mixins: [DragDropContext(HTML5Backend)],
  data() {
    return {
      over: {},
      dragging: {},
      expanded: [],
      search: null,
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
    },
    drag(dragging) {
      this.dragging = dragging;
    },
    enter(target) {
      // dragging used to be argument.check whether code is missing here
      this.expanded.push(target.id);
    },
    leave() {
      // dragging, target used to be arguments. check whether code is missing here
      this.over = null;
    },
    hover(dragging, target) {
      let parent = this.findParent(dragging.id, this.nodeTree);
      if (target.id !== parent.id) {
        this.over = { id: target.id, mode: "append" };
      }
    },
    drop(dragging, target) {
      let parent = this.findParent(dragging.id, this.nodeTree);
      if (dragging.id !== target.id && target.id !== parent.id) {
        let items = JSON.parse(JSON.stringify(this.nodeTree));
        // get dragging item (local copy!)
        let item = this.findItem(dragging.id, items);
        // remove from parent
        let draggingParent = this.findParent(item.id, items);
        let draggingIdx = draggingParent.children.findIndex(
          sibling => sibling === item
        );
        draggingParent.children.splice(draggingIdx, 1);

        // find target (local copy!)
        target = this.findItem(target.id, items);
        // add to target (local copy!)
        target.children.push(item);

        // sort
        if (this.config.options.ddAppendOnly) {
          target.children.sort((a, b) =>
            a["text"].localeCompare(b["text"], undefined, {
              sensitivity: "base"
            })
          );
        }

        // expand target
        if (this.expanded.indexOf(target.id) === -1) {
          this.expanded.push(target.id);
        }
      }
      this.over = null;
    },
    findParent(input, id) {
      for (let element of input) {
        if (element.id === id) {
          return element;
        } else if (element.children) {
          this.findParent(element.children, id);
        }
      }
    }
  },
  created() {
    this.$store.dispatch("fetchNodeTree").catch(e => console.log(e));
  },
  components: {
    DragDropSlot
  }
};
</script>

<style>
</style>
