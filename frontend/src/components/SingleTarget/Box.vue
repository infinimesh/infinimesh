<template>
  <div class="box" v-dragSource v-dropTarget v-bind:class="{dragging: isDragging, active: isActive, canDrop: !isActive && canDrop }">{{ name }}</div>
</template>

<script>
import { DragSource } from "vue-react-dnd";
import { DropTarget } from "vue-react-dnd";

export default {
  name: "Box",
  props: ["name"],
  mixins: [DragSource, DropTarget],

  data() {
    return {
      isDragging: false,
      isOver: false,
      canDrop: false
    };
  },

  dragSource: {
    type() {
      return "box";
    },
    specs: {
      beginDrag() {
        return {
          name: this.name
        };
      },

      endDrag(monitor) {
        const item = monitor.getItem();
        const dropResult = monitor.getDropResult();

        if (dropResult) {
          alert("You dropped " + item.name + " into " + dropResult.name + "!");
        }
      }
    },
    collect(connect, monitor) {
      this.isDragging = monitor.isDragging();
    }
  },
  computed: {
    isActive() {
      return this.canDrop && this.isOver;
    },

    text() {
      return this.isActive ? "Release to drop" : "Drag a box here";
    }
  },

  dropTarget: {
    type() {
      return "box";
    },
    specs: {
      drop() {
        return { name: "safdlök" };
      }
    },
    collect(connect, monitor) {
      this.isOver = monitor.isOver();
      this.canDrop = monitor.canDrop();
    }
  }
};
</script>

<style scoped>
.box {
  border: 1px solid gray;
  background-color: white;
  padding: 0.5rem 1rem;
  margin-right: 1.5rem;
  margin-bottom: 1.5rem;
  cursor: move;
  float: left;

  &.dragging {
    opacity: 0.4;
  }
}
</style>
