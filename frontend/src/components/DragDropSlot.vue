<template>
  <div v-dragSource v-dropTarget>
    <slot></slot>
  </div>
</template>

<script>
import { DragSource, DropTarget } from "vue-react-dnd";

export default {
  name: "DragDropSlot",
  mixins: [DragSource, DropTarget],
  props: {
    item: Object
  },
  data() {
    return {
      isOver: false
    };
  },
  methods: {
    emitDropTargetEvent(event, monitor) {
      const dragging = monitor.getItem();
      const target = this.item;
      this.$emit(event, dragging, target, monitor);
    }
  },
  dragSource: {
    type() {
      return "treeItem";
    },
    specs: {
      canDrag() {
        return this.item.allowDrag !== false && this.item.draggable !== false;
      },
      beginDrag() {
        this.$emit("drag", this.item);
        console.log("item dragged", JSON.stringify(this.item.id));
        return this.item;
      },
      endDrag() {
        this.$emit("drag", null);
      }
    }
  },
  dropTarget: {
    type() {
      return "treeItem";
    },
    specs: {
      canDrop(monitor) {
        const dragging = monitor.getItem();
        const target = this.item;
        return target.allowDrop !== false && dragging.id !== target.id;
      },
      hover(monitor) {
        if (monitor.canDrop() && monitor.isOver({ shallow: true })) {
          this.emitDropTargetEvent("hover", monitor);
        }
      },
      drop(monitor) {
        this.emitDropTargetEvent("drop", monitor);
      }
    },
    collect(connect, monitor) {
      const dragging = monitor.getItem();
      const target = this.item;

      if (dragging === null || target === null || dragging.id !== target.id) {
        let isOver = monitor.isOver({ shallow: true });

        if (isOver && !this.isOver) {
          this.emitDropTargetEvent("enter", monitor);
        } else if (!isOver && this.isOver) {
          this.emitDropTargetEvent("leave", monitor);
        }

        this.isOver = isOver;
      }
    }
  }
};
</script>




<!-- <template>
  <Container />
</template>

<script>
import Container from "./SingleTarget/Container.vue"

export default {
  components: {
    Container
  }
}
</script>

<style lang="css" scoped>
</style> -->
