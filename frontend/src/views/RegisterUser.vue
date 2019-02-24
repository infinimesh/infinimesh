<template>
  <div>
    <drag-drop-slot>
    </drag-drop-slot>
  </div>

</template>

<script>
import DragDropSlot from "../components/DragDropSlot.vue";

export default {
  components: {
    DragDropSlot
  }
};
</script>

<style lang="css" scoped>
</style>



<!-- <template>
  <v-container>
    <h1 class="mb-3">For now the tree drag and drop test page</h1>
    <v-card>
      <v-card-text>
        <v-treeview
          :items="items"
        >
        <template slot="label" slot-scope="{ item }">
          <drag-drop-slot :class="['tree-item', (over && over.id === item.id ? over.mode : '')]" :key="item.id" :item="item" @drag="drag" @enter="enter" @leave="leave" @hover="hover" @drop="drop">
            <v-icon>add</v-icon>
            <span>{{ item.name }}</span>
          </drag-drop-slot>
        </template>
        </v-treeview>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import DragDropSlot from "../components/DragDropSlot.vue";

export default {
  data() {
    return {
      over: {},
      items: [
        {
          id: 1,
          name: "Applications :",
          children: [
            { id: 2, name: "Calendar : app" },
            { id: 3, name: "Chrome : app" },
            { id: 4, name: "Webstorm : app" }
          ]
        },
        {
          id: 5,
          name: "Documents :",
          children: [
            {
              id: 6,
              name: "vuetify :",
              children: [
                {
                  id: 7,
                  name: "src :",
                  children: [
                    { id: 8, name: "index : ts" },
                    { id: 9, name: "bootstrap : ts" }
                  ]
                }
              ]
            },
            {
              id: 10,
              name: "material2 :",
              children: [
                {
                  id: 11,
                  name: "src :",
                  children: [
                    { id: 12, name: "v-btn : ts" },
                    { id: 13, name: "v-card : ts" },
                    { id: 14, name: "v-window : ts" }
                  ]
                }
              ]
            }
          ]
        },
        {
          id: 15,
          name: "Downloads :",
          children: [
            { id: 16, name: "October : pdf" },
            { id: 17, name: "November : pdf" },
            { id: 18, name: "Tutorial : html" }
          ]
        },
        {
          id: 19,
          name: "Videos :",
          children: [
            {
              id: 20,
              name: "Tutorials :",
              children: [
                { id: 21, name: "Basic layouts : mp4" },
                { id: 22, name: "Advanced techniques : mp4" },
                { id: 23, name: "All about app : dir" }
              ]
            },
            { id: 24, name: "Intro : mov" },
            { id: 25, name: "Conference introduction : avi" }
          ]
        }
      ]
    };
  },
  methods: {
    drag(dragging) {
      this.dragging = dragging;
    },
    enter(dragging, target) {
      this.expanded.push(target.id);
    },
    leave() {
      // leave had dragging and target as args
      this.over = null;
    },
    hover(dragging, target) {
      let parent = this.findParent(dragging.id, this.tree);
      if (target.id !== parent.id) {
        this.over = { id: target.id, mode: "append" };
      }
    },
    drop(dragging, target) {
      let parent = this.findParent(dragging.id, this.tree);
      if (dragging.id !== target.id && target.id !== parent.id) {
        let items = JSON.parse(JSON.stringify(this.tree));
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
    }
  },
  components: {
    DragDropSlot
  }
};
</script>

<style scoped>
</style> -->
